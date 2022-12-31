package metaevent

import (
	"fmt"

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
	unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error
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

func Parse(tick uint64, bytes ...byte) (any, error) {
	var delta lib.Delta
	var status uint8
	var eventType lib.MetaEventType
	var data []byte

	if v, remaining, err := vlq(bytes); err != nil {
		return nil, err
	} else if len(remaining) < 1 {
		return nil, fmt.Errorf("Invalid metaevent - missing status")
	} else if remaining[0] != 0xff {
		return nil, fmt.Errorf("Invalid metaevent status byte (%02X)", remaining[0])
	} else if len(remaining) < 2 {
		return nil, fmt.Errorf("Invalid metaevent - missing event type")
	} else if u, err := vlf(remaining[2:]); err != nil {
		return nil, err
	} else {
		delta = lib.Delta(v)
		status = remaining[0]
		eventType = lib.MetaEventType(remaining[1] & 0x7F)
		data = u
	}

	switch eventType {
	case lib.TypeSequenceNumber:
		return unmarshal[SequenceNumber](tick, delta, status, data, bytes...)

	case lib.TypeText:
		return unmarshal[Text](tick, delta, status, data, bytes...)

	case lib.TypeCopyright:
		return unmarshal[Copyright](tick, delta, status, data, bytes...)

	case lib.TypeTrackName:
		return unmarshal[TrackName](tick, delta, status, data, bytes...)

	case lib.TypeInstrumentName:
		return unmarshal[InstrumentName](tick, delta, status, data, bytes...)

	case lib.TypeLyric:
		return unmarshal[Lyric](tick, delta, status, data, bytes...)

	case lib.TypeMarker:
		return unmarshal[Marker](tick, delta, status, data, bytes...)

	case lib.TypeCuePoint:
		return unmarshal[CuePoint](tick, delta, status, data, bytes...)

	case lib.TypeProgramName:
		return unmarshal[ProgramName](tick, delta, status, data, bytes...)

	case lib.TypeDeviceName:
		return unmarshal[DeviceName](tick, delta, status, data, bytes...)

	case lib.TypeMIDIChannelPrefix:
		return unmarshal[MIDIChannelPrefix](tick, delta, status, data, bytes...)

	case lib.TypeMIDIPort:
		return unmarshal[MIDIPort](tick, delta, status, data, bytes...)

	case lib.TypeEndOfTrack:
		return unmarshal[EndOfTrack](tick, delta, status, data, bytes...)

	case lib.TypeTempo:
		return unmarshal[Tempo](tick, delta, status, data, bytes...)

	case lib.TypeSMPTEOffset:
		return unmarshal[SMPTEOffset](tick, delta, status, data, bytes...)

	case lib.TypeKeySignature:
		return unmarshal[KeySignature](tick, delta, status, data, bytes...)

	case lib.TypeTimeSignature:
		return unmarshal[TimeSignature](tick, delta, status, data, bytes...)

	case lib.TypeSequencerSpecificEvent:
		return unmarshal[SequencerSpecificEvent](tick, delta, status, data, bytes...)

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
	}](tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) (any, error) {
	p := P(new(E))
	if err := p.unmarshal(tick, delta, status, data, bytes...); err != nil {
		return nil, err
	} else {
		return *p, nil
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}

func equals[T lib.MetaEventType](b byte, t T) bool {
	return (b & 0x7f) == byte(t)
}

func vlq(bytes []byte) (uint32, []byte, error) {
	vlq := uint32(0)

	for i, b := range bytes {
		vlq <<= 7
		vlq += uint32(b & 0x7f)

		if b&0x80 == 0 {
			return vlq, bytes[i+1:], nil
		}
	}

	return 0, nil, fmt.Errorf("Invalid event 'delta'")
}

func vlf(bytes []byte) ([]byte, error) {
	if N, remaining, err := vlq(bytes); err != nil {
		return nil, err
	} else {
		return remaining[:N], nil
	}
}

func delta(bytes []byte) (lib.Delta, []byte, error) {
	v, remaining, err := vlq(bytes)

	return lib.Delta(v), remaining, err
}
