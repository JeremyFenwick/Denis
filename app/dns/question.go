package dns

import "encoding/binary"

type Question struct {
	Labels []string // Together this forms the domain name
	Type   uint16
	Class  uint16
}

func (question *Question) Encode() []byte {
	data := make([]byte, question.byteSize())
	index := 0

	// Add the labels
	for _, label := range question.Labels {
		// Add the label length
		data[index] = byte(len(label))
		index++
		// Add the label data
		copy(data[index:], label)
		index += len(label)
	}

	// Null terminator
	data[index] = 0x00
	index++

	// Add the type
	binary.BigEndian.PutUint16(data[index:], question.Type)
	index += 2

	// Add the class
	binary.BigEndian.PutUint16(data[index:], question.Class)
	index += 2 // Now at the end

	return data
}

func CodecraftersQuestion() *Question {
	return &Question{
		Labels: []string{"codecrafters", "io"},
		Type:   1,
		Class:  1,
	}
}

func (question *Question) byteSize() int {
	res := 0
	// Label + 1 for each label
	for _, label := range question.Labels {
		res += len(label) + 1
	}
	// Null terminator, type and class added at the end
	return res + 5
}
