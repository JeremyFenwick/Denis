package dns

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	// Decode the request header
	requestHeader, err := HeaderDecode(ctx.Data)
	if err != nil {
		ctx.Logger.Println("Error decoding header:", err)
		return
	}
	ctx.Logger.Printf("request header: %s", requestHeader)

	// Build the response header
	responseHeader := &Header{
		Id:      requestHeader.Id,
		Qr:      true,
		OpCode:  requestHeader.OpCode,
		Rd:      requestHeader.Rd,
		RCode:   getRCode(requestHeader.OpCode),
		QdCount: 1,
		AnCount: 1,
	}
	ctx.Logger.Printf("response header: %s", responseHeader)

	// Build the response message
	message := append(responseHeader.Encode(), CodecraftersQuestion().Encode()...)
	message = append(message, CodecraftersAnswer().Encode()...)

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
