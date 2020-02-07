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
	Type types.MetaEventType
}

func (e MetaEvent) String() string {
	return fmt.Sprintf("%s %s", e.Event, e.Type)
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

	if b, err := rr.ReadByte(); err != nil {
		return nil, err
	} else {
		event.Type = types.MetaEventType(b & 0x7F)
	}

	switch event.Type {
	case 0x00:
		return NewSequenceNumber(&event, rr)

	case 0x01:
		return NewText(&event, rr)

	case 0x02:
		return NewCopyright(&event, rr)

	case 0x03:
		return NewTrackName(&event, rr)

	case 0x04:
		return NewInstrumentName(&event, rr)

	case 0x05:
		return NewLyric(&event, rr)

	case 0x06:
		return NewMarker(&event, rr)

	case 0x07:
		return NewCuePoint(&event, rr)

	case 0x08:
		return NewProgramName(&event, rr)

	case 0x09:
		return NewDeviceName(&event, rr)

	case 0x20:
		return NewMIDIChannelPrefix(&event, rr)

	case 0x21:
		return NewMIDIPort(&event, rr)

	case 0x2f:
		return NewEndOfTrack(&event, rr)

	case 0x51:
		return NewTempo(&event, rr)

	case 0x54:
		return NewSMPTEOffset(&event, rr)

	case 0x58:
		return NewTimeSignature(&event, rr)

	case 0x59:
		return NewKeySignature(ctx, &event, rr)

	case 0x7f:
		return NewSequencerSpecificEvent(ctx, &event, rr)
	}

	return nil, fmt.Errorf("Unrecognised META event: %02X", byte(event.Type))
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
