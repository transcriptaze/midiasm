package metaevent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/types"
)

func Parse(ctx *context.Context, r io.ByteReader, status types.Status) (interface{}, error) {
	if status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%v): expected 'FF'", status)
	}

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	eventType := types.MetaEventType(b & 0x7F)

	switch eventType {
	case 0x00:
		return NewSequenceNumber(data)

	case 0x01:
		return NewText(data)

	case 0x02:
		return NewCopyright(data)

	case 0x03:
		return NewTrackName(data)

	case 0x04:
		return NewInstrumentName(data)

	case 0x05:
		return NewLyric(data)

	case 0x06:
		return NewMarker(data)

	case 0x07:
		return NewCuePoint(data)

	case 0x08:
		return NewProgramName(data)

	case 0x09:
		return NewDeviceName(data)

	case 0x20:
		return NewMIDIChannelPrefix(data)

	case 0x21:
		return NewMIDIPort(data)

	case 0x51:
		return NewTempo(data)

	case 0x54:
		return NewSMPTEOffset(data)

	case 0x58:
		return NewTimeSignature(data)

	case 0x59:
		return NewKeySignature(ctx, data)

	case 0x2f:
		return NewEndOfTrack(data)

	case 0x7f:
		return NewSequencerSpecificEvent(data)
	}

	return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
}
