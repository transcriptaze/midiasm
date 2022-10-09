package events

import (
	"io"

	"github.com/transcriptaze/midiasm/midi/types"
)

type Event struct {
	tick  uint64
	Delta types.Delta
	Event any
	Bytes types.Hex `json:"-"`
}

func NewEvent(tick uint64, delta uint32, evt any, bytes []byte) *Event {
	return &Event{
		tick:  tick,
		Delta: types.Delta(delta),
		Event: evt,
		Bytes: bytes,
	}
}

func (e Event) Tick() uint64 {
	return uint64(e.tick)
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
