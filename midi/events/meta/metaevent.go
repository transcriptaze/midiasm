package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TMetaEvent interface {
	SequenceNumber |
		Text |
		Copyright |
		TrackName |
		InstrumentName |
		Lyric |
		Marker |
		CuePoint |
		ProgramName |
		DeviceName |
		MIDIChannelPrefix |
		MIDIPort |
		EndOfTrack |
		Tempo |
		SMPTEOffset |
		KeySignature |
		TimeSignature |
		SequencerSpecificEvent
}

type IMetaEvent interface {
	unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error
}

type event struct {
	tick   uint64
	delta  lib.Delta
	bytes  []byte
	tag    lib.Tag
	Status lib.Status
	Type   lib.MetaEventType
}

func (e event) Tick() uint64 {
	return e.tick
}

func (e event) Delta() uint32 {
	return uint32(e.delta)
}

func (e event) Bytes() []byte {
	return e.bytes
}

func (e event) Tag() string {
	return fmt.Sprintf("%v", e.tag)
}

func Parse(ctx *context.Context, tick uint64, delta lib.Delta, status byte, b byte, data []byte, bytes ...byte) (any, error) {
	eventType := lib.MetaEventType(b & 0x7F)

	switch eventType {
	case lib.TypeSequenceNumber:
		if e, err := UnmarshalSequenceNumber(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeText:
		if e, err := UnmarshalText(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeCopyright:
		if e, err := UnmarshalCopyright(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeTrackName:
		if e, err := UnmarshalTrackName(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeInstrumentName:
		if e, err := UnmarshalInstrumentName(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeLyric:
		if e, err := UnmarshalLyric(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMarker:
		if e, err := UnmarshalMarker(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeCuePoint:
		if e, err := UnmarshalCuePoint(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeProgramName:
		if e, err := UnmarshalProgramName(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeDeviceName:
		if e, err := UnmarshalDeviceName(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMIDIChannelPrefix:
		if e, err := UnmarshalMIDIChannelPrefix(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMIDIPort:
		if e, err := UnmarshalMIDIPort(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeEndOfTrack:
		if e, err := UnmarshalEndOfTrack(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeTempo:
		if e, err := UnmarshalTempo(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeSMPTEOffset:
		if e, err := UnmarshalSMPTEOffset(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeKeySignature:
		if e, err := UnmarshalKeySignature(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeTimeSignature:
		if e, err := UnmarshalTimeSignature(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeSequencerSpecificEvent:
		if e, err := UnmarshalSequencerSpecificEvent(ctx, tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	default:
		return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
	}
}

// Ref. https://stackoverflow.com/questions/71444847/go-with-generics-type-t-is-pointer-to-type-parameter-not-type-parameter
// Ref. https://stackoverflow.com/questions/69573113/how-can-i-instantiate-a-non-nil-pointer-of-type-argument-with-generic-go/69575720#69575720
func unmarshal[
	E TMetaEvent,
	P interface {
		*E
		IMetaEvent
	}](ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	p := P(new(E))
	if err := p.unmarshal(ctx, tick, delta, status, data, bytes...); err != nil {
		return nil, err
	} else {
		return *p, nil
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}
