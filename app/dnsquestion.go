package main

type DNSQuestion struct {
	Labels []string
	Type   uint16
	Class  uint16
}

func (question *DNSQuestion) Encode() []byte {
	data := make([]byte, 0)
	for i := 0; i < len(question.Labels); i++ {
		// Add the label length
		data = append(data, byte(len(question.Labels[i])))
		// Add the label data
		data = append(data, []byte(question.Labels[i])...)

	}
	data = append(data, 0x00) // Null terminator
	// Add the type and class
	data = append(data, byte(question.Type>>8), byte(question.Type&0xFF))
	data = append(data, byte(question.Class>>8), byte(question.Class&0xFF))
	return data
}

func CodecraftersQuestion() *DNSQuestion {
	return &DNSQuestion{
		Labels: []string{"codecrafters", "io"},
		Type:   1,
		Class:  1,
	}
}
