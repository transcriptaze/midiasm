package transpose

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Transpose struct {
	Writer io.Writer
}

func (t *Transpose) Execute(b []byte) ([]byte, error) {
	chunks := [][]byte{}
	r := bufio.NewReader(bytes.NewReader(b))

	// ... MThd
	if chunk, err := readChunk(r); err != nil {
		return nil, err
	} else if chunk == nil {
		return nil, fmt.Errorf("missing MThd chunk")
	} else if string(chunk[0:4]) != "MThd" {
		return nil, fmt.Errorf("invalid MThd chunk")
	} else {
		tag := chunk[0:4]
		length := chunk[4:8]

		debugf("%-4v %v %v", len(chunk), string(tag), binary.BigEndian.Uint32(length))

		chunks = append(chunks, chunk)
	}

	// ... tracks
	for {
		chunk, err := readChunk(r)

		if err != nil && err == io.EOF {
			return flatten(chunks), nil
		} else if err != nil {
			return nil, err
		} else if chunk == nil {
			return nil, fmt.Errorf("invalid chunk (%v)", chunk)
		}

		// ... MTrk?
		if string(chunk[0:4]) == "MTrk" {
			if list, err := transpose(chunk); err != nil {
				return nil, err
			} else {
				chunks = append(chunks, list...)
			}

			continue
		}

		// ... default
		chunks = append(chunks, chunk)
	}
}

func transpose(track []byte) ([][]byte, error) {
	tag := track[0:4]
	length := track[4:8]
	chunks := [][]byte{
		tag,
		length,
	}

	debugf("%-4v %v %v", len(track), string(tag), binary.BigEndian.Uint32(length))

	r := bufio.NewReader(bytes.NewReader(track[8:]))

	if evt, err := event(r); err != nil {
		return nil, err
	} else {
		chunks = append(chunks, evt)
	}

	if evt, err := event(r); err != nil {
		return nil, err
	} else {
		chunks = append(chunks, evt)
	}

	if remaining, err := io.ReadAll(r); err != nil {
		return nil, err
	} else {
		fmt.Printf("remaining: %v\n", len(remaining))
		chunks = append(chunks, remaining)
	}

	return chunks, nil
}

type reader struct {
	r      *bufio.Reader
	buffer *bytes.Buffer
}

func (r reader) ReadByte() (byte, error) {
	if b, err := r.r.ReadByte(); err != nil {
		return b, err
	} else {
		return b, r.buffer.WriteByte(b)
	}
}

func (r reader) Peek(n int) ([]byte, error) {
	return r.r.Peek(n)
}

func event(rr *bufio.Reader) ([]byte, error) {
	var buffer bytes.Buffer

	r := reader{
		r:      bufio.NewReader(rr),
		buffer: &buffer,
	}

	if _, err := VLQ(r); err != nil {
		return nil, err
	} else if b, err := r.Peek(1); err != nil {
		return nil, err
	} else {
		switch b[0] {
		case 0xff: // ... meta event
			if status, err := r.ReadByte(); err != nil {
				return nil, err
			} else if event, err := r.ReadByte(); err != nil {
				return nil, err
			} else if data, err := VLF(r); err != nil {
				return nil, err
			} else {
				debugf("%-4v      META  status:%02X event:%02X data:%v", buffer.Len(), status, event, data)
				return buffer.Bytes(), nil
			}
		}
	}

	// // ... SysEx event
	// if b == 0xf0 || b == 0xf7 {
	// 	ctx.RunningStatus = 0x00

	// 	rr.ReadByte()

	// 	e, err := sysex.Parse(rr, types.Status(b), ctx)

	// 	return &events.Event{
	// 		Tick:  types.Tick(tick + delta),
	// 		Delta: types.Delta(delta),
	// 		Bytes: buffer.Bytes(),
	// 		Event: e,
	// 	}, err
	// }

	// // ... MIDI event
	// if b < 0x80 && ctx.RunningStatus < 0x80 {
	// 	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", b&0xF0)
	// }

	// status := types.Status(b)

	// if b < 0x80 {
	// 	status = ctx.RunningStatus
	// } else {
	// 	rr.ReadByte()
	// }

	// ctx.RunningStatus = status

	// e, err := midievent.Parse(rr, status, ctx)

	// return &events.Event{
	// 	Tick:  types.Tick(tick + delta),
	// 	Delta: types.Delta(delta),
	// 	Bytes: buffer.Bytes(),
	// 	Event: e,
	// }, err

	return buffer.Bytes(), nil
}

func readChunk(r *bufio.Reader) ([]byte, error) {
	peek, err := r.Peek(8)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(peek[4:8])
	bytes := make([]byte, length+8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func VLQ(r io.ByteReader) (uint32, error) {
	vlq := uint32(0)

	for {
		b, err := r.ReadByte()
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

func VLF(r io.ByteReader) ([]byte, error) {
	N, err := VLQ(r)
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

func flatten(chunks [][]byte) []byte {
	var b bytes.Buffer

	for _, chunk := range chunks {
		b.Write(chunk)
	}

	return b.Bytes()
}

func debugf(format string, args ...any) {
	f := fmt.Sprintf("%-6v  %v\n", "DEBUG", format)

	fmt.Printf(f, args...)
}
