package midi

import ()

type vlq struct {
	v uint32
}

func (v vlq) MarshalBinary() ([]byte, error) {
	buffer := []byte{0x00, 0x80, 0x80, 0x80, 0x00}
	b := v.v

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
