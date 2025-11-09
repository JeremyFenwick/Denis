package main

import (
	"fmt"

	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	var err error

	// Based on the configuration, generate a response and send it to the client
	if ctx.ForwardingOn {
		err = forwardingResponse(ctx)
	} else {
		err = simpleResponse(ctx)
	}
	if err != nil {
		ctx.Logger.Println("Error generating response:", err)
		return
	}
}

// Response formats
func forwardingResponse(ctx *udp_server.PacketContext) error {
	ctx.Logger.Println("Using forwarding mode")

	// Parse the request header
	requestHeader, questions, err := parseRequest(ctx)
	if err != nil {
		return err
	}

	// Forward requests and collect all answers
	allAnswers, rcode, err := forwardRequests(requestHeader, questions, ctx)
	if err != nil {
		ctx.Logger.Println("Error forwarding requests:", err)
		return err
	}

	ctx.Logger.Printf("Received %d answers from forwarder", len(allAnswers))

	// Build the final response with all answers
	responseHeader := &dns.Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Aa:      false,
		Tc:      false,
		Rd:      requestHeader.Rd,
		Ra:      false,
		Z:       0,
		RCode:   rcode, // Use RCode from forwarder
		QdCount: uint16(len(questions)),
		AnCount: uint16(len(allAnswers)),
		NsCount: 0,
		ArCount: 0,
	}

	// Encode the complete response
	message, err := encodeResponse(responseHeader, questions, allAnswers)
	if err != nil {
		return err
	}

	// Send the response
	err = ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
		return err
	}

	return nil
}

func simpleResponse(ctx *udp_server.PacketContext) error {
	// Parse the request
	requestHeader, questions, err := parseRequest(ctx)
	if err != nil {
		return err
	}

	// Build the response
	responseHeader := buildResponseHeader(requestHeader, len(questions))
	answers := generateAnswers(questions)

	// Encode the response
	message, err := encodeResponse(responseHeader, questions, answers)
	if err != nil {
		return err
	}

	// Send the response
	err = ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}

	return nil
}

// Helper functions

func forwardRequests(header *dns.Header, questions []*dns.Question, ctx *udp_server.PacketContext) ([]*dns.Answer, uint8, error) {
	allAnswers := make([]*dns.Answer, 0)
	var finalRCode uint8 = 0

	// Generate a request for each question
	for i := 0; i < len(questions); i++ {
		ctx.Logger.Printf("Forwarding question %d/%d", i+1, len(questions))

		// Generate a single request for the forwarder
		newHeader := *header // Shallow copy
		newHeader.QdCount = 1
		request := newHeader.Encode()
		qEncoded, err := questions[i].Encode()
		if err != nil {
			return nil, 0, fmt.Errorf("error encoding question: %w", err)
		}
		request = append(request, qEncoded...)

		ctx.Logger.Printf("Sending %d bytes to forwarder", len(request))

		// Send the request to the forwarder, receive the response
		response, err := ctx.ForwardAndReceive(request)
		if err != nil {
			return nil, 0, fmt.Errorf("error forwarding request: %w", err)
		}

		ctx.Logger.Printf("Received %d bytes from forwarder", len(response))

		// Parse the answer from the response
		answers, rcode, err := parseAnswersFromResponse(response)
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing response: %w", err)
		}

		ctx.Logger.Printf("Parsed %d answers with RCode=%d", len(answers), rcode)

		// Use the RCode from the first response (or you could use the last non-zero one)
		if i == 0 || rcode != 0 {
			finalRCode = rcode
		}

		// Collect all answers
		allAnswers = append(allAnswers, answers...)
	}

	return allAnswers, finalRCode, nil
}

func parseAnswersFromResponse(response []byte) ([]*dns.Answer, uint8, error) {
	// Decode the response header
	responseHeader, err := dns.HeaderDecode(response)
	if err != nil {
		return nil, 0, fmt.Errorf("error decoding response header: %w", err)
	}

	// Skip the header (12 bytes)
	index := 12

	// Skip the question section
	for i := 0; i < int(responseHeader.QdCount); i++ {
		_, bytesUsed, err := dns.QuestionDecode(response, index)
		if err != nil {
			return nil, 0, fmt.Errorf("error decoding question: %w", err)
		}
		index += bytesUsed
	}

	// Decode the answer section
	answers := make([]*dns.Answer, responseHeader.AnCount)
	for i := 0; i < int(responseHeader.AnCount); i++ {
		answer, bytesUsed, err := dns.AnswerDecode(response, index)
		if err != nil {
			return nil, 0, fmt.Errorf("error decoding answer: %w", err)
		}
		answers[i] = answer
		index += bytesUsed
	}

	return answers, responseHeader.RCode, nil
}

func parseRequest(ctx *udp_server.PacketContext) (*dns.Header, []*dns.Question, error) {
	// Decode the request header
	requestHeader, err := dns.HeaderDecode(ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding header:", err)
		return nil, nil, err
	}
	ctx.Logger.Printf("request: %s", requestHeader)

	// Decode the request questions
	questions, err := decodeQuestions(int(requestHeader.QdCount), ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding questions:", err)
		return nil, nil, err
	}
	for _, question := range questions {
		ctx.Logger.Printf("request: %s", question)
	}

	return requestHeader, questions, nil
}

func buildResponseHeader(requestHeader *dns.Header, questionCount int) *dns.Header {
	return &dns.Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Rd:      requestHeader.Rd,
		RCode:   getRCode(requestHeader.OpCode),
		QdCount: uint16(questionCount),
		AnCount: uint16(questionCount),
	}
}

func encodeResponse(header *dns.Header, questions []*dns.Question, answers []*dns.Answer) ([]byte, error) {
	response := header.Encode()

	// Add questions
	for _, question := range questions {
		encodedQuestion, err := question.Encode()
		if err != nil {
			return nil, fmt.Errorf("error encoding question: %w", err)
		}
		response = append(response, encodedQuestion...)
	}

	// Add answers
	for _, answer := range answers {
		encodedAnswer, err := answer.Encode()
		if err != nil {
			return nil, fmt.Errorf("error encoding answer: %w", err)
		}
		response = append(response, encodedAnswer...)
	}

	return response, nil
}

func decodeQuestions(count int, data []byte) ([]*dns.Question, error) {
	index := 12 // Skip the header
	result := make([]*dns.Question, count)

	for i := 0; i < count; i++ {
		question, bytesUsed, err := dns.QuestionDecode(data, index)
		if err != nil {
			return nil, err
		}
		result[i] = question
		index += bytesUsed
	}

	return result, nil
}

func generateAnswers(questions []*dns.Question) []*dns.Answer {
	answers := make([]*dns.Answer, len(questions))
	for i, question := range questions {
		answers[i] = &dns.Answer{
			Label:  question.Label,
			Type:   1,
			Class:  1,
			TTL:    60,
			Length: 4,
			Data:   []byte{8, 8, 8, 8},
		}
	}
	return answers
}

func getRCode(opCode uint8) uint8 {
	if opCode == 0 {
		return 0
	}
	return 4
}
