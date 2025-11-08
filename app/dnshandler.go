package main

import "github.com/codecrafters-io/dns-server-starter-go/app/server"

func DnsHandler(ctx *server.PacketContext) {
	responseHeader := DnsHeader{
		Id:      1234,
		Qr:      true,
		QdCount: 1,
	}
	message := append(responseHeader.Encode(), CodecraftersQuestion().Encode()...)
	err := ctx.Send(message, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}
