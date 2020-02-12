package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type MetaEvent struct {
	events.Event
}

type reader struct {
	rdr   io.ByteReader
	event *MetaEvent
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.Bytes = append(r.event.Bytes, b)
	}

	return b, err
}

func Parse(e events.Event, r io.ByteReader, ctx *context.Context) (interface{}, error) {
	if e.Status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%02x): expected 'FF'", e.Status)
	}

	ctx.ClearRunningStatus()

	event := MetaEvent{
		Event: e,
	}

	rr := reader{r, &event}

	b, err := rr.ReadByte()
	if err != nil {
		return nil, err
	}

	eventType := types.MetaEventType(b & 0x7F)

	switch eventType {
	case 0x00:
		return NewSequenceNumber(&event, eventType, rr)

	case 0x01:
		return NewText(&event, eventType, rr)

	case 0x02:
		return NewCopyright(&event, eventType, rr)

	case 0x03:
		return NewTrackName(&event, eventType, rr)

	case 0x04:
		return NewInstrumentName(&event, eventType, rr)

	case 0x05:
		return NewLyric(&event, eventType, rr)

	case 0x06:
		return NewMarker(&event, eventType, rr)

	case 0x07:
		return NewCuePoint(&event, eventType, rr)

	case 0x08:
		return NewProgramName(&event, eventType, rr)

	case 0x09:
		return NewDeviceName(&event, eventType, rr)

	case 0x20:
		return NewMIDIChannelPrefix(&event, eventType, rr)

	case 0x21:
		return NewMIDIPort(&event, eventType, rr)

	case 0x2f:
		return NewEndOfTrack(&event, eventType, rr)

	case 0x51:
		return NewTempo(&event, eventType, rr)

	case 0x54:
		return NewSMPTEOffset(&event, eventType, rr)

	case 0x58:
		return NewTimeSignature(&event, eventType, rr)

	case 0x59:
		return NewKeySignature(ctx, &event, eventType, rr)

	case 0x7f:
		return NewSequencerSpecificEvent(ctx, &event, eventType, rr)
	}

	return nil, fmt.Errorf("Unrecognised META event: %02X", byte(eventType))
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
