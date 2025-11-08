package dns

import (
	"encoding/binary"
	"fmt"
)

type Answer struct {
	Label  Label // Together this forms the domain name
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   []byte
}

func (answer *Answer) Encode() ([]byte, error) {
	data := make([]byte, answer.byteSize())
	index := 0

	// Add the labels
	bytesUsed, err := answer.Label.encode(data)
	if err != nil {
		return nil, err
	}
	index += bytesUsed
	// Add the null terminator
	data[index] = 0x00
	index++

	// Add the type
	binary.BigEndian.PutUint16(data[index:], answer.Type)
	index += 2

	// Add the class
	binary.BigEndian.PutUint16(data[index:], answer.Class)
	index += 2

	// Add the TTL
	binary.BigEndian.PutUint32(data[index:], answer.TTL)
	index += 4

	// Add the length
	binary.BigEndian.PutUint16(data[index:], answer.Length)
	index += 2

	// Add the answer data
	copy(data[index:], answer.Data)

	return data, nil
}

func (answer *Answer) String() string {
	return fmt.Sprintf(
		"Answer{Labels: %s, Type: %d, Class: %d, TTL: %d, Length: %d, Data: %s}",
		answer.Label, answer.Type, answer.Class, answer.TTL, answer.Length, answer.Data)
}

func CodecraftersAnswer() *Answer {
	return &Answer{
		Label:  NewLabel([]string{"codecrafters", "io"}),
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   []byte{8, 8, 8, 8},
	}
}

func (answer *Answer) byteSize() int {
	res := 0
	// Get the labels
	res += answer.Label.ByteSize()
	// Add the type, class, TTL and length
	res += 10
	// Add the data
	res += len(answer.Data)
	return res
}
