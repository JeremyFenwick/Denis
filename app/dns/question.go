package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Question struct {
	Label Label // domain name
	Type  uint16
	Class uint16
}

func (question *Question) Encode() ([]byte, error) {
	data := make([]byte, question.byteSize())
	index := 0

	// Encode the label
	consumed, err := question.Label.encode(data)
	if err != nil {
		return nil, err
	}
	index += consumed

	// Null terminator
	data[index] = 0x00
	index++

	// Add type and class
	binary.BigEndian.PutUint16(data[index:], question.Type)
	binary.BigEndian.PutUint16(data[index+2:], question.Class)

	return data, nil
}

func QuestionDecode(data []byte, start int) (*Question, int, error) {
	label, bytesUsed, err := LabelDecode(data, start)
	if err != nil {
		return nil, 0, err
	}

	index := start + bytesUsed
	if len(data[index:]) < 4 {
		return nil, 0, errors.New("provided question data is too small")
	}

	question := &Question{
		Label: label,
		Type:  binary.BigEndian.Uint16(data[index:]),
		Class: binary.BigEndian.Uint16(data[index+2:]),
	}
	return question, bytesUsed + 4, nil
}

func (question *Question) String() string {
	return fmt.Sprintf("Question{Labels: %s, Type: %d, Class: %d}", question.Label, question.Type, question.Class)
}

func (question *Question) byteSize() int {
	res := 0
	// Labels size
	res += question.Label.ByteSize()
	// type and class added at the end
	return res + 4
}
