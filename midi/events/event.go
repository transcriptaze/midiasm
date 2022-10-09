package events

import (
	"io"

	"github.com/transcriptaze/midiasm/midi/types"
)

type IEvent interface {
	Tick() uint64
	Delta() uint32
	Bytes() []byte
}

type Event struct {
	tick  uint64
	delta uint32
	Event any
	bytes types.Hex `json:"-"`
}

func NewEvent(tick uint64, delta uint32, evt any, bytes []byte) *Event {
	return &Event{
		tick:  tick,
		delta: delta,
		Event: evt,
		bytes: bytes,
	}
}

func (e Event) Tick() uint64 {
	if v, ok := e.Event.(IEvent); ok {
		return v.Tick()
	}

	return e.tick
}

func (e Event) Delta() uint32 {
	if v, ok := e.Event.(IEvent); ok {
		return v.Delta()
	}

	return e.delta
}

func (e Event) Bytes() types.Hex {
	if v, ok := e.Event.(IEvent); ok {
		return v.Bytes()
	}

	return e.bytes
}

// FIXME: temporary hack while reworking Event as embedded struct
func (e *Event) Mutate(offset int, b byte) {
	e.bytes[len(e.bytes)+offset] = b
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
