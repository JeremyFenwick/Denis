package main

import (
	"fmt"

	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	var err error

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

func forwardingResponse(ctx *udp_server.PacketContext) error {
	// Parse the request
	requestHeader, questions, err := parseRequest(ctx)
	if err != nil {
		return err
	}

	// Collect all forwarded responses
	var allAnswers []*dns.Answer
	var rcode uint8

	for i, question := range questions {
		// Build a single-question request
		singleHeader := *requestHeader
		singleHeader.QdCount = 1
		request := singleHeader.Encode()

		questionBytes, err := question.Encode()
		if err != nil {
			return err
		}
		request = append(request, questionBytes...)

		// Forward and receive
		response, err := ctx.ForwardAndReceive(request)
		if err != nil {
			return err
		}

		// Decode the response header
		respHeader, err := dns.HeaderDecode(response)
		if err != nil {
			return err
		}

		// Extract RCode from first response
		if i == 0 {
			rcode = respHeader.RCode
		}

		// Parse answers from this response
		answers, err := extractAnswers(response, respHeader)
		if err != nil {
			return err
		}

		allAnswers = append(allAnswers, answers...)
	}

	// Build final response with original ID
	responseHeader := &dns.Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Rd:      requestHeader.Rd,
		RCode:   rcode,
		QdCount: uint16(len(questions)),
		AnCount: uint16(len(allAnswers)),
	}

	message, err := encodeResponse(responseHeader, questions, allAnswers)
	if err != nil {
		return err
	}

	return ctx.Send(message, ctx.Address)
}

func simpleResponse(ctx *udp_server.PacketContext) error {
	requestHeader, questions, err := parseRequest(ctx)
	if err != nil {
		return err
	}

	responseHeader := buildResponseHeader(requestHeader, len(questions))
	answers := generateAnswers(questions)

	message, err := encodeResponse(responseHeader, questions, answers)
	if err != nil {
		return err
	}

	return ctx.Send(message, ctx.Address)
}

func parseRequest(ctx *udp_server.PacketContext) (*dns.Header, []*dns.Question, error) {
	header, err := dns.HeaderDecode(ctx.Data)
	if err != nil {
		return nil, nil, err
	}

	questions, err := decodeQuestions(int(header.QdCount), ctx.Data)
	if err != nil {
		return nil, nil, err
	}

	return header, questions, nil
}

func extractAnswers(response []byte, header *dns.Header) ([]*dns.Answer, error) {
	index := 12 // Skip header

	// Skip questions
	for i := 0; i < int(header.QdCount); i++ {
		_, bytesUsed, err := dns.QuestionDecode(response, index)
		if err != nil {
			return nil, err
		}
		index += bytesUsed
	}

	// Read answers
	answers := make([]*dns.Answer, header.AnCount)
	for i := 0; i < int(header.AnCount); i++ {
		answer, bytesUsed, err := dns.AnswerDecode(response, index)
		if err != nil {
			return nil, err
		}
		answers[i] = answer
		index += bytesUsed
	}

	return answers, nil
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

	for _, question := range questions {
		encoded, err := question.Encode()
		if err != nil {
			return nil, fmt.Errorf("error encoding question: %w", err)
		}
		response = append(response, encoded...)
	}

	for _, answer := range answers {
		encoded, err := answer.Encode()
		if err != nil {
			return nil, fmt.Errorf("error encoding answer: %w", err)
		}
		response = append(response, encoded...)
	}

	return response, nil
}

func decodeQuestions(count int, data []byte) ([]*dns.Question, error) {
	index := 12
	questions := make([]*dns.Question, count)

	for i := 0; i < count; i++ {
		question, bytesUsed, err := dns.QuestionDecode(data, index)
		if err != nil {
			return nil, err
		}
		questions[i] = question
		index += bytesUsed
	}

	return questions, nil
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
