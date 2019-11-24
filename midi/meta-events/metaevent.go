package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type MetaEvent struct {
	event.Event
	eventType byte
	bytes     []byte
}

type reader struct {
	rdr   io.ByteReader
	event *MetaEvent
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.bytes = append(r.event.bytes, b)
	}

	return b, err
}

func Parse(e event.Event, x []byte, r io.ByteReader) (event.IEvent, error) {
	if e.Status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%02x): expected 'FF'", e.Status)
	}

	event := MetaEvent{
		Event: e,
		bytes: append(make([]byte, 0), x...),
	}

	rr := reader{r, &event}

	if b, err := rr.ReadByte(); err != nil {
		return nil, err
	} else {
		event.eventType = b & 0x7F
	}

	switch event.eventType {
	case 0x03:
		return NewTrackName(&event, rr)

	case 0x2f:
		return NewEndOfTrack(&event, rr)

	case 0x51:
		return NewTempo(&event, rr)

	case 0x58:
		return NewTimeSignature(&event, rr)

	case 0x59:
		return NewKeySignature(&event, rr)
	}

	return nil, fmt.Errorf("Unrecognised META event: %02X", event.eventType)
}

func read(r io.ByteReader) ([]byte, error) {
	N, err := vlq(r)
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

func vlq(r io.ByteReader) (uint32, error) {
	l := uint32(0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		l <<= 7
		l += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return l, nil
}
