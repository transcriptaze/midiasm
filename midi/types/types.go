package types

type VLQ uint32

func (v VLQ) MarshalBinary() ([]byte, error) {
	return vlq2bin(uint32(v)), nil
}

type VLF []byte

func (v VLF) MarshalBinary() ([]byte, error) {
	return vlf2bin(v), nil
}

func vlq2bin(v uint32) []byte {
	buffer := []byte{0x00, 0x80, 0x80, 0x80, 0x00}

	for i := 4; i > 0; i-- {
		buffer[i] |= byte(v & 0x7f)
		if v >>= 7; v == 0 {
			return buffer[i:]
		}
	}

	buffer[1] |= 0x80
	buffer[0] = byte(v & 0x7f)

	return buffer
}

func vlf2bin(v []byte) []byte {
	len := uint32(len(v))

	return append(vlq2bin(len), v...)
}
