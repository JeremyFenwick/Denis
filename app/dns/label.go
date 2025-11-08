package dns

import (
	"errors"
	"fmt"
)

type Label []string

func LabelDecode(data []byte, start int) (Label, int, error) {
	label, consumed, err := decoderLoop(start, data)
	if err != nil {
		return nil, 0, err
	}
	return label, consumed, nil
}

func (label *Label) encode(data []byte) (int, error) {
	if len(data) < label.ByteSize() {
		return 0, errors.New("provided data is too small")
	}
	offset := 0
	for _, entry := range *label {
		// Add the label length
		data[offset] = byte(len(entry))
		offset++
		// Add the label data
		copy(data[offset:], entry)
		offset += len(entry)
	}
	return offset, nil
}

func NewLabel(entries []string) Label {
	return Label(entries)
}

func decoderLoop(start int, data []byte) (Label, int, error) {
	index := start
	content := make([]string, 0)

	for {
		if index >= len(data) {
			return nil, 0, errors.New("provided labels data is too small")
		}
		// Break on the null terminator
		if data[index] == 0x00 {
			index++
			break
		}
		// Follow a pointer
		if data[index]&0xC0 == 0xC0 {
			offset := getLocation(data[index], data[index+1])
			toAdd, _, err := decoderLoop(offset, data)
			fmt.Printf("Following pointer at index %d to offset %d (data len: %d)\n", index, offset, len(data))
			if err != nil {
				return nil, 0, err
			}
			content = append(content, toAdd...)
			index += 2
			break
		}
		// Otherwise we just add the label
		// Get the length
		length := int(data[index])
		index++
		content = append(content, string(data[index:index+length]))
		index += length
	}

	return content, index - start, nil
}
func getLocation(first byte, second byte) int {
	return (int(first&0x3F) << 8) | int(second)
}

func (label *Label) ByteSize() int {
	res := 0
	// Label + 1 for each label
	for _, entry := range *label {
		res += len(entry) + 1
	}
	// Add the null terminator
	return res + 1
}
