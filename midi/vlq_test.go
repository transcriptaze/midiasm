package midi

import (
	"reflect"
	"testing"
)

func TestVLQMarshalBinary(t *testing.T) {
	tests := []struct {
		v        uint32
		expected []byte
	}{
		{0, []byte{0x00}},
		{480, []byte{0x83, 0x60}},
	}

	for _, test := range tests {
		v := vlq{test.v}
		if encoded, err := v.MarshalBinary(); err != nil {
			t.Fatalf("Error encoding VLQ value %v (%v)", test.v, err)
		} else if !reflect.DeepEqual(encoded, test.expected) {
			t.Errorf("Incorrectly encoded VLQ value %v - expected:%x, got:%x", test.v, test.expected, encoded)
		}
	}
}
