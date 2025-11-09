package udp_server

import (
	"errors"
	"log"
	"net"
	"time"
)

type Config struct {
	Address        string
	Logger         *log.Logger
	Handler        HandlerFunc
	ForwardAddress string
}

// UdpServer - handler function, then listens for incoming packets
type UdpServer struct {
	conn           *net.UDPConn
	bufferSize     int
	logger         *log.Logger
	handler        HandlerFunc
	forwardAddress *net.UDPAddr
}

type HandlerFunc func(ctx *PacketContext)

func NewUdpServer(config Config) (*UdpServer, error) {
	udpAddress, err := net.ResolveUDPAddr("udp", config.Address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return nil, err
	}

	server := &UdpServer{
		conn:       conn,
		bufferSize: 1024,
		logger:     config.Logger,
		handler:    config.Handler,
	}

	// If a forward address is provided, set it
	if config.ForwardAddress != "" {
		var forwardAddress *net.UDPAddr
		forwardAddress, err = net.ResolveUDPAddr("udp", config.ForwardAddress)
		if err != nil {
			return nil, err
		}
		server.forwardAddress = forwardAddress
	}

	return server, nil
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
			Data:              buffer[:bytesRead],
			Address:           address,
			Logger:            server.logger,
			Send:              server.Send,
			ForwardingOn:      server.forwardAddress != nil,
			ForwardAndReceive: server.ForwardAndReceive,
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

func (server *UdpServer) ForwardAndReceive(packet []byte) ([]byte, error) {
	if server.forwardAddress == nil {
		return nil, errors.New("no forward address configured")
	}

	// Create a new UDP connection for this request
	// (so we can receive the response on a specific port)
	conn, err := net.DialUDP("udp", nil, server.forwardAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Set a timeout so we don't wait forever
	err = conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send the request
	_, err = conn.Write(packet)
	if err != nil {
		return nil, err
	}

	// Receive the response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
