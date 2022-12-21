package events

import (
	"encoding"
	"fmt"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func (e *Event) UnmarshalBinary(bytes []byte) error {
	if _, remaining, err := vlq(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else {
		status := remaining[0]

		switch {
		case status == 0xff && equals(remaining[1], lib.TypeSequenceNumber):
			return unmarshalBinary[metaevent.SequenceNumber](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeText):
			return unmarshalBinary[metaevent.Text](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeCopyright):
			return unmarshalBinary[metaevent.Copyright](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeTrackName):
			return unmarshalBinary[metaevent.TrackName](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeInstrumentName):
			return unmarshalBinary[metaevent.InstrumentName](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeLyric):
			return unmarshalBinary[metaevent.Lyric](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeMarker):
			return unmarshalBinary[metaevent.Marker](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeCuePoint):
			return unmarshalBinary[metaevent.CuePoint](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeProgramName):
			return unmarshalBinary[metaevent.ProgramName](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeDeviceName):
			return unmarshalBinary[metaevent.DeviceName](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeMIDIChannelPrefix):
			return unmarshalBinary[metaevent.MIDIChannelPrefix](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeMIDIPort):
			return unmarshalBinary[metaevent.MIDIPort](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeTempo):
			return unmarshalBinary[metaevent.Tempo](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeTimeSignature):
			return unmarshalBinary[metaevent.TimeSignature](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeKeySignature):
			return unmarshalBinary[metaevent.KeySignature](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeSMPTEOffset):
			return unmarshalBinary[metaevent.SMPTEOffset](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeEndOfTrack):
			return unmarshalBinary[metaevent.EndOfTrack](e, bytes)

		case status == 0xff && equals(remaining[1], lib.TypeSequencerSpecificEvent):
			return unmarshalBinary[metaevent.SequencerSpecificEvent](e, bytes)

		case equals(status, lib.TypeNoteOff):
			return unmarshalBinary[midievent.NoteOff](e, bytes)

		case equals(status, lib.TypeNoteOn):
			return unmarshalBinary[midievent.NoteOn](e, bytes)

		case equals(status, lib.TypePolyphonicPressure):
			return unmarshalBinary[midievent.PolyphonicPressure](e, bytes)

		case equals(status, lib.TypeController):
			return unmarshalBinary[midievent.Controller](e, bytes)

		case equals(status, lib.TypeProgramChange):
			return unmarshalBinary[midievent.ProgramChange](e, bytes)

		case equals(status, lib.TypeChannelPressure):
			return unmarshalBinary[midievent.ChannelPressure](e, bytes)

		case equals(status, lib.TypePitchBend):
			return unmarshalBinary[midievent.PitchBend](e, bytes)

		case equals(status, lib.TypeSysExMessage):
			return unmarshalBinary[sysex.SysExMessage](e, bytes)

		case equals(status, lib.TypeSysExContinuationMessage):
			return unmarshalBinary[sysex.SysExContinuationMessage](e, bytes)

		// case "SysExEscape":
		//     return unmarshalJSON[sysex.SysExEscapeMessage](e, t.Event)

		default:
			return fmt.Errorf("Unrecognised event (%02X)", status)
		}
	}
}

func unmarshalBinary[
	E TEvent,
	P interface {
		*E
		encoding.BinaryUnmarshaler
	}](e *Event, bytes []byte) (err error) {
	p := P(new(E))

	if err = p.UnmarshalBinary(bytes); err != nil {
		return err
	} else {
		e.Event = *p
	}

	return
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

func equals[T lib.TEventType](b byte, t T) bool {
	return t.Equals(b)
}

// func equals[T lib.MidiEventType](b byte, t T) bool {
// 	return (b & 0x70) == byte(t)
// }
