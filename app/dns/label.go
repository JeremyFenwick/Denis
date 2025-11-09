package dns

import (
	"errors"
	"fmt"
)

type Label []string

func LabelDecode(data []byte, start int) (Label, int, error) {
	visited := make(map[int]bool)
	label, consumed, err := decoderLoop(start, data, visited)
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

// decoderLoop now takes a visited map to detect cycles
func decoderLoop(start int, data []byte, visited map[int]bool) (Label, int, error) {
	// Check for cycles - if we've already visited this offset, we have a circular reference
	if visited[start] {
		return nil, 0, fmt.Errorf("circular pointer detected at offset %d", start)
	}
	visited[start] = true

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
			if index+1 >= len(data) {
				return nil, 0, errors.New("pointer extends beyond data")
			}
			offset := getLocation(data[index], data[index+1])

			// Validate pointer offset
			if offset >= len(data) {
				return nil, 0, fmt.Errorf("pointer offset %d exceeds data length %d", offset, len(data))
			}
			// DNS pointers should typically point backward, but we'll just rely on cycle detection

			toAdd, _, err := decoderLoop(offset, data, visited)
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

		// Validate length doesn't exceed remaining data
		if index+length > len(data) {
			return nil, 0, fmt.Errorf("label length %d exceeds remaining data at index %d", length, index)
		}

		// Validate label length (DNS labels max 63 bytes)
		if length > 63 {
			return nil, 0, fmt.Errorf("label length %d exceeds maximum of 63 bytes", length)
		}

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
