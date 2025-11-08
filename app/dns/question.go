package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

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

func QuestionDecode(data []byte) (*Question, error) {
	// Get the labels
	var labels []string
	index := 0
	for {
		if index >= len(data) {
			return nil, errors.New("provided question data is too small")
		}
		labelLength := int(data[index])
		index++
		// Break on the null terminator
		if labelLength == 0 {
			break
		}
		labels = append(labels, string(data[index:index+labelLength]))
		index += labelLength
	}

	// If we don't have enough data, return an error
	if len(data[index:]) < 4 {
		return nil, errors.New("provided question data is too small")
	}

	// Get the type and class
	question := &Question{
		Labels: labels,
		Type:   binary.BigEndian.Uint16(data[index:]),
		Class:  binary.BigEndian.Uint16(data[index+2:]),
	}
	return question, nil
}

func CodecraftersQuestion() *Question {
	return &Question{
		Labels: []string{"codecrafters", "io"},
		Type:   1,
		Class:  1,
	}
}

func (question *Question) String() string {
	return fmt.Sprintf("Question{Labels: %s, Type: %d, Class: %d}", question.Labels, question.Type, question.Class)
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
