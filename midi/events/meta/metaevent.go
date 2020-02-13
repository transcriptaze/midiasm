package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

func Parse(ctx *context.Context, r io.ByteReader, status types.Status) (interface{}, error) {
	if status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%v): expected 'FF'", status)
	}

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	eventType := types.MetaEventType(b & 0x7F)

	switch eventType {
	case 0x00:
		return NewSequenceNumber(r)

	case 0x01:
		return NewText(r)

	case 0x02:
		return NewCopyright(r)

	case 0x03:
		return NewTrackName(r)

	case 0x04:
		return NewInstrumentName(r)

	case 0x05:
		return NewLyric(r)

	case 0x06:
		return NewMarker(r)

	case 0x07:
		return NewCuePoint(r)

	case 0x08:
		return NewProgramName(r)

	case 0x09:
		return NewDeviceName(r)

	case 0x20:
		return NewMIDIChannelPrefix(r)

	case 0x21:
		return NewMIDIPort(r)

	case 0x2f:
		return NewEndOfTrack(r)

	case 0x51:
		return NewTempo(r)

	case 0x54:
		return NewSMPTEOffset(r)

	case 0x58:
		return NewTimeSignature(r)

	case 0x59:
		return NewKeySignature(ctx, r)

	case 0x7f:
		return NewSequencerSpecificEvent(ctx, r)
	}

	return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
}
