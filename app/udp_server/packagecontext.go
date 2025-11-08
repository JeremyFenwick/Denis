package udp_server

import (
	"log"
	"net"
)

type PacketContext struct {
	Data    []byte
	Address *net.UDPAddr
	Logger  *log.Logger
	Send    SendFunc
}

type SendFunc func(packet []byte, address *net.UDPAddr) error
