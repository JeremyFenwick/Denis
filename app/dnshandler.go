package main

import "github.com/codecrafters-io/dns-server-starter-go/app/server"

func DnsHandler(ctx *server.PacketContext) {
	responseHeader := DnsHeader{
		Id: 1234,
		Qr: true,
	}
	err := ctx.Send(responseHeader.Encode(), ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}
