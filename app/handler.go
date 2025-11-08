package main

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	// Decode the request header
	requestHeader, err := dns.HeaderDecode(ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding header:", err)
		return
	}
	ctx.Logger.Printf("request: %s", requestHeader)

	// Decode the request question
	question, err := dns.QuestionDecode(ctx.Data[12:])
	if err != nil {
		ctx.Logger.Println("Error decoding question:", err)
		return
	}
	ctx.Logger.Printf("request: %s", question)

	// Build the response header
	responseHeader := &dns.Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Rd:      requestHeader.Rd,
		RCode:   getRCode(requestHeader.OpCode),
		QdCount: 1,
		AnCount: 1,
	}
	ctx.Logger.Printf("response header: %s", responseHeader)

	// Build the response question
	responseQuestion := &dns.Question{
		Labels: question.Labels,
		Type:   1,
		Class:  1,
	}
	ctx.Logger.Printf("response question: %s", responseQuestion)

	// Build the response answer
	responseAnswer := &dns.Answer{
		Labels: question.Labels,
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   []byte{8, 8, 8, 8},
	}
	ctx.Logger.Printf("response answer: %s", responseAnswer)

	// Build the response message
	message := append(responseHeader.Encode(), responseQuestion.Encode()...)
	message = append(message, responseAnswer.Encode()...)

	err = ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}

func getRCode(opCode uint8) uint8 {
	if opCode == 0 {
		return 0
	}
	return 4
}
