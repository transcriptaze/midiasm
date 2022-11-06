package types

import (
	"fmt"
	"strconv"
)

type Delta uint32

func (d Delta) String() string {
	return fmt.Sprintf("%d", uint32(d))
}

func ParseDelta(s string) (Delta, error) {
	if delta, err := strconv.ParseUint(s, 10, 32); err != nil {
		return 0, err
	} else {
		return Delta(delta), nil
	}
}

func (d Delta) MarshalBinary() ([]byte, error) {
	buffer := []byte{0x00, 0x80, 0x80, 0x80, 0x00}
	b := d

	for i := 4; i > 0; i-- {
		buffer[i] |= byte(b & 0x7f)
		if b >>= 7; b == 0 {
			return buffer[i:], nil
		}
	}

	buffer[1] |= 0x80
	buffer[0] = byte(b & 0x7f)

	return buffer, nil
}
