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
		return NewSequenceNumber(r, status, eventType)

	case 0x01:
		return NewText(r, status, eventType)

	case 0x02:
		return NewCopyright(r, status, eventType)

	case 0x03:
		return NewTrackName(r, status, eventType)

	case 0x04:
		return NewInstrumentName(r, status, eventType)

	case 0x05:
		return NewLyric(r, status, eventType)

	case 0x06:
		return NewMarker(r, status, eventType)

	case 0x07:
		return NewCuePoint(r, status, eventType)

	case 0x08:
		return NewProgramName(r, status, eventType)

	case 0x09:
		return NewDeviceName(r, status, eventType)

	case 0x20:
		return NewMIDIChannelPrefix(r, status, eventType)

	case 0x21:
		return NewMIDIPort(r, status, eventType)

	case 0x2f:
		return NewEndOfTrack(r, status, eventType)

	case 0x51:
		return NewTempo(r, status, eventType)

	case 0x54:
		return NewSMPTEOffset(r, status, eventType)

	case 0x58:
		return NewTimeSignature(r, status, eventType)

	case 0x59:
		return NewKeySignature(ctx, r, status, eventType)

	case 0x7f:
		return NewSequencerSpecificEvent(ctx, r, status, eventType)
	}

	return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
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
