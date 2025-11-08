package main

import "github.com/codecrafters-io/dns-server-starter-go/app/server"

func DnsHandler(ctx *server.PacketContext) {
	var response []byte
	err := ctx.Send(response, ctx.Address)
	if err != nil {
		ctx.Logger.Println("Error sending response:", err)
	}
}
