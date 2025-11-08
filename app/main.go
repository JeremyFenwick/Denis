package main

import (
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/app/server"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	logger := log.New(os.Stdout, "[DNS Server] ", log.LstdFlags)
	udpServer, err := server.NewUdpServer("127.0.0.1:2053", logger, DnsHandler)
	if err != nil {
		logger.Fatal("> Error creating UDP server:", err)
	}
	udpServer.Listen()
}
