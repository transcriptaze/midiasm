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
			t.Errorf("Incorrectlry parsed hex string - expected:%v, got:%v", test.hex, h)
		}
	}
}
