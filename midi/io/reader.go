package IO

import (
	"bufio"
	"bytes"
)

type Reader struct {
	reader *bufio.Reader
	buffer *bytes.Buffer
}

func NewReader(r *bufio.Reader) Reader {
	return Reader{
		reader: r,
		buffer: &bytes.Buffer{},
	}
}

func BytesReader(b []byte) Reader {
	r := bufio.NewReader(bytes.NewBuffer(b))

	return Reader{
		reader: r,
		buffer: &bytes.Buffer{},
	}
}

func TestReader(buffered []byte, b []byte) Reader {
	r := bufio.NewReader(bytes.NewBuffer(b))

	return Reader{
		reader: r,
		buffer: bytes.NewBuffer(buffered),
	}
}

func (r Reader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return b, err
	}

	return b, r.buffer.WriteByte(b)
}

func (r Reader) Read(buffer []byte) (int, error) {
	N, err := r.reader.Read(buffer)
	if err != nil {
		return N, err
	}

	r.buffer.Write(buffer)

	return N, nil
}

func (r Reader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}

func (r Reader) Bytes() []byte {
	return r.buffer.Bytes()
}
