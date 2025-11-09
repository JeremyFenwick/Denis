package main

import (
	"flag"
	"log"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/app/udp_server"
)

func main() {
	// Parse command-line flags
	resolver := flag.String("resolver", "", "DNS resolver address (ip:port)")
	flag.Parse()

	// Setup the logger
	logger := log.New(os.Stdout, "[DNS Server] ", log.LstdFlags)

	// UDP server config
	config := udp_server.Config{
		Address: "127.0.0.1:2053",
		Logger:  logger,
		Handler: Handler,
	}

	// Configure forwarding if resolver flag is provided
	if *resolver != "" {
		config.ForwardAddress = *resolver
		logger.Printf("Forwarding enabled to resolver: %s", *resolver)
	} else {
		logger.Println("Running in simple mode (no forwarding)")
	}

	// Create and start the server
	udpServer, err := udp_server.NewUdpServer(config)
	if err != nil {
		logger.Fatal("Error creating UDP server: ", err)
	}
	udpServer.Listen()
}
