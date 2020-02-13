package sysex

import (
	"bufio"
)

type reader struct {
	r *bufio.Reader
}

func (r reader) ReadByte() (byte, error) {
	return r.r.ReadByte()
}

func (r reader) Peek(n int) ([]byte, error) {
	return r.r.Peek(n)
}

func (r reader) ReadVLQ() ([]byte, error) {
	N, err := r.VLQ()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, N)

	for i := 0; i < int(N); i++ {
		if b, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			bytes[i] = b
		}
	}

	return bytes, nil
}

func (r reader) VLQ() (uint32, error) {
	vlq := uint32(0)

	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}

		vlq <<= 7
		vlq += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return vlq, nil
}
