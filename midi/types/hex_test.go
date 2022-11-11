package types

import (
	"reflect"
	"testing"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		s   string
		hex []byte
	}{
		{"00 09 01", []byte{0x00, 0x09, 0x01}},
	}

	for _, test := range tests {
		if h, err := ParseHex(test.s); err != nil {
			t.Errorf("Error parsing valid hex string (%v)", err)
		} else if !reflect.DeepEqual(h, Hex(test.hex)) {
			t.Errorf("Incorrectly parsed hex string - expected:%v, got:%v", test.hex, h)
		}
	}
}

func TestHexMarshalBinary(t *testing.T) {
	tests := []struct {
		hex      []byte
		expected []byte
	}{
		{[]byte{0x7e, 0x00, 0x09, 0x01}, []byte{0x04, 0x7e, 0x00, 0x09, 0x01}},
	}

	for _, test := range tests {
		h := Hex(test.hex)
		if bytes, err := h.MarshalBinary(); err != nil {
			t.Errorf("Error marshalling hex data (%v)", err)
		} else if !reflect.DeepEqual(bytes, test.expected) {
			t.Errorf("Incorrectly marshalled hex data- expected:%v, got:%v", test.expected, bytes)
		}
	}
}
