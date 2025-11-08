package udp_server

import (
	"log"
	"net"
)

type UdpServer struct {
	conn       *net.UDPConn
	bufferSize int
	logger     *log.Logger
	handler    HandlerFunc
}

type HandlerFunc func(ctx *PacketContext)

func NewUdpServer(address string, logger *log.Logger, handler HandlerFunc) (*UdpServer, error) {
	udpAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return nil, err
	}

	return &UdpServer{
		conn:       conn,
		bufferSize: 1024,
		logger:     logger,
		handler:    handler,
	}, nil
}

func (server *UdpServer) Close() {
	_ = server.conn.Close()
}

func (server *UdpServer) Listen() {
	defer server.Close()

	buffer := make([]byte, server.bufferSize)

	for {
		// Get the packet
		bytesRead, address, err := server.conn.ReadFromUDP(buffer)
		if err != nil {
			server.logger.Println("Error receiving data:", err)
			break
		}
		server.logger.Printf("Recieved %d bytes from %s", bytesRead, address)

		// Handle the packet in a goroutine
		ctx := &PacketContext{
			Data:    buffer[:bytesRead],
			Address: address,
			Logger:  server.logger,
			Send:    server.Send,
		}
		go server.handler(ctx)
	}
}

func (server *UdpServer) Send(packet []byte, address *net.UDPAddr) error {
	bytesWritten, err := server.conn.WriteToUDP(packet, address)
	if err != nil {
		server.logger.Printf("Sent %d bytes to %s", bytesWritten, address)
	}
	return err
}
