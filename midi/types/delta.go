package types

import (
	"fmt"
	"regexp"
	"strconv"
)

type Delta uint32

func (d Delta) String() string {
	return fmt.Sprintf("%d", uint32(d))
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

func (d *Delta) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile("delta:([0-9]+)")

	if match := re.FindStringSubmatch(string(bytes)); match == nil || len(match) != 2 {
		return fmt.Errorf("missing delta field")
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		*d = Delta(delta)
	}

	return nil
}
