package main

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	message, err := generateResponse(ctx)

	// Send the response
	err = ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}

func generateResponse(ctx *udp_server.PacketContext) ([]byte, error) {
	// Decode the request header
	requestHeader, err := dns.HeaderDecode(ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding header:", err)
		return nil, err
	}
	ctx.Logger.Printf("request: %s", requestHeader)
	// Decode the request question/s
	questions, err := decodeQuestions(int(requestHeader.QdCount), ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding question:", err)
		return nil, err
	}
	for _, question := range questions {
		ctx.Logger.Printf("request: %s", question)
	}
	// Generate the answers
	answers := generateAnswers(questions)
	// Build the response header
	responseHeader := &dns.Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Rd:      requestHeader.Rd,
		RCode:   getRCode(requestHeader.OpCode),
		QdCount: uint16(len(questions)),
		AnCount: uint16(len(answers)),
	}
	response := responseHeader.Encode()
	// Add the questions (un-compressed)
	for _, question := range questions {
		encodedQuestion, err := question.Encode()
		if err != nil {
			ctx.Logger.Println("Error encoding question:", err)
			return nil, err
		}
		response = append(response, encodedQuestion...)
	}
	// Add the answers
	for _, answer := range answers {
		encodedAnswer, err := answer.Encode()
		if err != nil {
			ctx.Logger.Println("Error encoding answer:", err)
			return nil, err
		}
		response = append(response, encodedAnswer...)
	}
	return response, nil
}

func decodeQuestions(count int, data []byte) ([]*dns.Question, error) {
	index := 12
	result := make([]*dns.Question, count)

	for i := 0; i < count; i++ {
		question, bytesUsed, err := dns.QuestionDecode(data, index) // ⭐ Pass full data + index
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
		answer := &dns.Answer{
			Label:  question.Label,
			Type:   1,
			Class:  1,
			TTL:    60,
			Length: 4,
			Data:   []byte{8, 8, 8, 8},
		}
		answers[i] = answer
	}
	return answers
}

func getRCode(opCode uint8) uint8 {
	if opCode == 0 {
		return 0
	}
	return 4
}
