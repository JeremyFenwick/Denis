package dns_test

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
)

func TestAnswerEncodeDecode(t *testing.T) {
	lbl := dns.NewLabel([]string{"google", "com"})

	a := &dns.Answer{
		Label:  lbl,
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   []byte{1, 2, 3, 4},
	}

	enc, err := a.Encode()
	if err != nil {
		t.Fatal(err)
	}

	dec, used, err := dns.AnswerDecode(enc, 0)
	if err != nil {
		t.Fatal(err)
	}

	if used != len(enc) {
		t.Fatalf("expected %d bytes used, got %d", len(enc), used)
	}

	if !bytes.Equal(dec.Data, a.Data) {
		t.Fatal("data mismatch")
	}
}

func TestHeaderEncodeDecode(t *testing.T) {
	h := &dns.Header{
		Id: 1234, Qr: true, OpCode: 1, Aa: true, Rd: true,
		RCode: 3, QdCount: 1, AnCount: 2, NsCount: 3, ArCount: 4,
	}

	enc := h.Encode()
	dec, err := dns.HeaderDecode(enc)
	if err != nil {
		t.Fatal(err)
	}

	if dec.Id != h.Id || dec.OpCode != h.OpCode || dec.RCode != h.RCode || dec.AnCount != h.AnCount {
		t.Fatal("header fields mismatch")
	}
}

func TestQuestionEncodeDecode(t *testing.T) {
	lbl := dns.NewLabel([]string{"google", "com"})
	q := &dns.Question{Label: lbl, Type: 1, Class: 1}

	enc, err := q.Encode()
	if err != nil {
		t.Fatal(err)
	}

	dec, used, err := dns.QuestionDecode(enc, 0)
	if err != nil {
		t.Fatal(err)
	}

	if used != len(enc) {
		t.Fatal("used mismatch")
	}

	if dec.Type != q.Type || dec.Class != q.Class {
		t.Fatal("value mismatch")
	}
}

func TestHeaderDecodeTooSmall(t *testing.T) {
	_, err := dns.HeaderDecode([]byte{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for too small header")
	}
}

func TestLabelCircularPointer(t *testing.T) {
	//  C0 00 means pointer to 0 offset -> infinite loop if not handled
	data := []byte{0xC0, 0x00}
	_, _, err := dns.LabelDecode(data, 0)
	if err == nil {
		t.Fatal("expected circular pointer error")
	}
}

func TestLabelInvalidLength(t *testing.T) {
	// 70 length means label length 112 which exceeds 63 limit
	data := []byte{70, 'a', 'b', 'c'}
	_, _, err := dns.LabelDecode(data, 0)
	if err == nil {
		t.Fatal("expected invalid label length error")
	}
}

func TestAnswerDecodeInsufficientData(t *testing.T) {
	lbl := dns.NewLabel([]string{"google", "com"})
	a := &dns.Answer{
		Label: lbl,
		Type:  1, Class: 1, TTL: 60,
		Length: 4,
		Data:   []byte{1, 2, 3, 4},
	}
	enc, _ := a.Encode()

	// chop last byte off -> should error
	bad := enc[:len(enc)-1]

	_, _, err := dns.AnswerDecode(bad, 0)
	if err == nil {
		t.Fatal("expected insufficient data error")
	}
}
