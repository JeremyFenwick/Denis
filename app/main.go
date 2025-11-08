package main

import (
	"log"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func main() {
	logger := log.New(os.Stdout, "[DNS Server] ", log.LstdFlags)
	udpServer, err := udp_server.NewUdpServer("127.0.0.1:2053", logger, Handler)
	if err != nil {
		logger.Fatal("> Error creating UDP udp_server:", err)
	}
	udpServer.Listen()
}
