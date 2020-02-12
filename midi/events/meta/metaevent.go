package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

func Parse(event *events.Event, r io.ByteReader, ctx *context.Context) (interface{}, error) {
	if event.Status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%02x): expected 'FF'", event.Status)
	}

	ctx.ClearRunningStatus()

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	eventType := types.MetaEventType(b & 0x7F)

	switch eventType {
	case 0x00:
		return NewSequenceNumber(event, eventType, r)

	case 0x01:
		return NewText(event, eventType, r)

	case 0x02:
		return NewCopyright(event, eventType, r)

	case 0x03:
		return NewTrackName(event, eventType, r)

	case 0x04:
		return NewInstrumentName(event, eventType, r)

	case 0x05:
		return NewLyric(event, eventType, r)

	case 0x06:
		return NewMarker(event, eventType, r)

	case 0x07:
		return NewCuePoint(event, eventType, r)

	case 0x08:
		return NewProgramName(event, eventType, r)

	case 0x09:
		return NewDeviceName(event, eventType, r)

	case 0x20:
		return NewMIDIChannelPrefix(event, eventType, r)

	case 0x21:
		return NewMIDIPort(event, eventType, r)

	case 0x2f:
		return NewEndOfTrack(event, eventType, r)

	case 0x51:
		return NewTempo(event, eventType, r)

	case 0x54:
		return NewSMPTEOffset(event, eventType, r)

	case 0x58:
		return NewTimeSignature(event, eventType, r)

	case 0x59:
		return NewKeySignature(ctx, event, eventType, r)

	case 0x7f:
		return NewSequencerSpecificEvent(ctx, event, eventType, r)
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
