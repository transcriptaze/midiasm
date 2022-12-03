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

type TMetaEventX interface {
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
		KeySignature |
		SequencerSpecificEvent

	MarshalJSON() ([]byte, error)
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
		if e, err := UnmarshalSequenceNumber(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeText:
		if e, err := UnmarshalText(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeCopyright:
		if e, err := UnmarshalCopyright(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeTrackName:
		if e, err := UnmarshalTrackName(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeInstrumentName:
		if e, err := UnmarshalInstrumentName(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeLyric:
		if e, err := UnmarshalLyric(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMarker:
		if e, err := UnmarshalMarker(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeCuePoint:
		if e, err := UnmarshalCuePoint(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeProgramName:
		if e, err := UnmarshalProgramName(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeDeviceName:
		if e, err := UnmarshalDeviceName(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMIDIChannelPrefix:
		if e, err := UnmarshalMIDIChannelPrefix(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeMIDIPort:
		if e, err := UnmarshalMIDIPort(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeEndOfTrack:
		if e, err := UnmarshalEndOfTrack(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeTempo:
		if e, err := UnmarshalTempo(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeSMPTEOffset:
		if e, err := UnmarshalSMPTEOffset(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeKeySignature:
		if e, err := UnmarshalKeySignature(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			if ctx != nil {
				if e.Accidentals < 0 {
					ctx.UseFlats()
				} else {
					ctx.UseSharps()
				}
			}

			return *e, err
		}

	case lib.TypeTimeSignature:
		if e, err := UnmarshalTimeSignature(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case lib.TypeSequencerSpecificEvent:
		if e, err := UnmarshalSequencerSpecificEvent(tick, delta, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	default:
		return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}
