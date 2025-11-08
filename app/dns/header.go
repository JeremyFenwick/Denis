package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Header struct {
	Id      uint16 // Packet ID
	Qr      bool   // Query/Response Indicator
	OpCode  uint8  // Operation Code
	Aa      bool   // Authoritative Answer
	Tc      bool   // Truncation
	Rd      bool   // Recursion Desired
	Ra      bool   // Recursion Available
	Z       uint8  // Reserved
	RCode   uint8  // Response Code
	QdCount uint16 // Question Count
	AnCount uint16 // Answer Record Count
	NsCount uint16 // Authority Record Count
	ArCount uint16 // Additional Record Count
}

func (header *Header) Encode() []byte {
	data := make([]byte, 12)
	// ID section
	binary.BigEndian.PutUint16(data[0:2], header.Id)
	// Flag section
	binary.BigEndian.PutUint16(data[2:4], header.getFlags())
	// Question section
	binary.BigEndian.PutUint16(data[4:6], header.QdCount)
	binary.BigEndian.PutUint16(data[6:8], header.AnCount)
	binary.BigEndian.PutUint16(data[8:10], header.NsCount)
	binary.BigEndian.PutUint16(data[10:12], header.ArCount)
	return data
}

func HeaderDecode(data []byte) (*Header, error) {
	if len(data) < 12 {
		return nil, errors.New("provided header is too small")
	}
	header := &Header{}
	// ID section
	header.Id = binary.BigEndian.Uint16(data[0:2])
	// Flag section
	flags := binary.BigEndian.Uint16(data[2:4])
	header.Qr = (flags & (1 << 15)) != 0
	header.OpCode = uint8((flags >> 11) & 0xF)
	header.Aa = (flags & (1 << 10)) != 0
	header.Tc = (flags & (1 << 9)) != 0
	header.Rd = (flags & (1 << 8)) != 0
	header.Ra = (flags & (1 << 7)) != 0
	header.Z = uint8((flags >> 4) & 0x7)
	header.RCode = uint8(flags & 0xF)
	// Question section
	header.QdCount = binary.BigEndian.Uint16(data[4:6])
	header.AnCount = binary.BigEndian.Uint16(data[6:8])
	header.NsCount = binary.BigEndian.Uint16(data[8:10])
	header.ArCount = binary.BigEndian.Uint16(data[10:12])

	return header, nil
}

// Helper functions
func (header *Header) getFlags() uint16 {
	var flags uint16
	// Qr
	if header.Qr {
		flags |= 1 << 15
	}
	// OpCode
	flags |= uint16(header.OpCode&0xF) << 11 // Mask to 4 bits
	// Aa
	if header.Aa {
		flags |= 1 << 10
	}
	// Tc
	if header.Tc {
		flags |= 1 << 9
	}
	// Rd
	if header.Rd {
		flags |= 1 << 8
	}
	// Ra
	if header.Ra {
		flags |= 1 << 7
	}
	// Z (reserved)
	flags |= uint16(header.Z&0x7) << 4 // Mask to 3 bits
	// RCode
	flags |= uint16(header.RCode & 0xF) // Mask to 4 bits
	return flags
}

func (header *Header) String() string {
	return fmt.Sprintf(
		"Header{ID: %d, QR: %t, OpCode: %d, AA: %t, TC: %t, RD: %t, RA: %t, Z: %d, RCode: %d, QDCount: %d, ANCount: %d, NSCount: %d, ARCount: %d}",
		header.Id, header.Qr, header.OpCode, header.Aa, header.Tc, header.Rd, header.Ra, header.Z, header.RCode, header.QdCount, header.AnCount, header.NsCount, header.ArCount,
	)
}
