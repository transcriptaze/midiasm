package midi

import (
	"bufio"
	"bytes"
)

type reader struct {
	reader *bufio.Reader
	buffer *bytes.Buffer
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return b, err
	}

	return b, r.buffer.WriteByte(b)
}

func (r reader) Read(buffer []byte) (int, error) {
	N, err := r.reader.Read(buffer)
	if err != nil {
		return N, err
	}

	r.buffer.Write(buffer)

	return N, nil
}

func (r reader) peek() (byte, error) {
	if bytes, err := r.reader.Peek(1); err != nil {
		return 0, err
	} else {
		return bytes[0], nil
	}
}

func (r reader) Bytes() []byte {
	return r.buffer.Bytes()
}
