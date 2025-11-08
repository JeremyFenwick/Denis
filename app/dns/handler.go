package dns

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func Handler(ctx *udp_server.PacketContext) {
	responseHeader := Header{
		Id:      1234,
		Qr:      true,
		QdCount: 1,
		AnCount: 1,
	}
	// Build the response message
	message := append(responseHeader.Encode(), CodecraftersQuestion().Encode()...)
	message = append(message, CodecraftersAnswer().Encode()...)

	err := ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}
