package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type MetaEvent struct {
	event.Event
	eventType byte
	length    uint32
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
	if e.Status != 0xff {
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
		event.eventType = b & 0x7f
	}

	l, err := vlq(rr)
	if err != nil {
		return nil, err
	}
	event.length = l

	ix := len(event.bytes)
	for i := 0; i < int(l); i++ {
		if _, err := rr.ReadByte(); err != nil {
			return nil, err
		}
	}

	switch event.eventType {
	case 0x03:
		return NewTrackName(event, event.bytes[ix:])

	case 0x2f:
		return NewEndOfTrack(event, event.bytes[ix:])

	case 0x51:
		return NewTempo(event, event.bytes[ix:])

	case 0x58:
		return NewTimeSignature(event, event.bytes[ix:])

	case 0x59:
		return NewKeySignature(event, event.bytes[ix:])
	}

	return nil, fmt.Errorf("Unrecognised META event: %02X", event.eventType)
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
