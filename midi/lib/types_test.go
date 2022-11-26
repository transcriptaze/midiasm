package lib

import (
	"reflect"
	"testing"
)

func TestVlq2Bin(t *testing.T) {
	tests := []struct {
		vlq      uint32
		expected []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{480, []byte{0x83, 0x60}},
	}

	for _, test := range tests {
		bytes := vlq2bin(test.vlq)
		if !reflect.DeepEqual(bytes, test.expected) {
			t.Errorf("Incorrectly marshalled VLQ value (%v) - expected:%v, got:%v", test.vlq, test.expected, bytes)
		}
	}
}
