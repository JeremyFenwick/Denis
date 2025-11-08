package dns

import "encoding/binary"

type Answer struct {
	Labels []string // Together this forms the domain name
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   []byte
}

func (answer *Answer) Encode() []byte {
	data := make([]byte, answer.byteSize())
	index := 0
	// Add the labels
	for _, label := range answer.Labels {
		// Add the label length
		data[index] = byte(len(label))
		index++
		// Add the label data
		copy(data[index:], label)
		index += len(label)
	}

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

	return data
}

func CodecraftersAnswer() *Answer {
	return &Answer{
		Labels: []string{"codecrafters", "io"},
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
	for _, label := range answer.Labels {
		res += len(label) + 1
	}
	// Add the null terminator
	res++
	// Add the type, class, TTL and length
	res += 10
	// Add the data
	res += len(answer.Data)
	return res
}
