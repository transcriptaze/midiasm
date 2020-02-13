package midi

import (
	"bufio"
	"bytes"
)

type reader struct {
	rdr    *bufio.Reader
	buffer *bytes.Buffer
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err != nil {
		return b, err
	}

	return b, r.buffer.WriteByte(b)
}

func (r reader) Peek(n int) ([]byte, error) {
	return r.rdr.Peek(n)
}
