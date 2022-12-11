package events

import (
	"encoding"
	"fmt"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/lib"
	// "github.com/transcriptaze/midiasm/midi/events/midi"
	// "github.com/transcriptaze/midiasm/midi/events/sysex"
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

		// case "Text":
		//     return unmarshalJSON[metaevent.Text](e, t.Event)

		// case "Copyright":
		//     return unmarshalJSON[metaevent.Copyright](e, t.Event)

		// case "TrackName":
		//     return unmarshalJSON[metaevent.TrackName](e, t.Event)

		// case "InstrumentName":
		//     return unmarshalJSON[metaevent.InstrumentName](e, t.Event)

		// case "Lyric":
		//     return unmarshalJSON[metaevent.Lyric](e, t.Event)

		// case "Marker":
		//     return unmarshalJSON[metaevent.Marker](e, t.Event)

		// case "CuePoint":
		//     return unmarshalJSON[metaevent.CuePoint](e, t.Event)

		// case "ProgramName":
		//     return unmarshalJSON[metaevent.ProgramName](e, t.Event)

		// case "DeviceName":
		//     return unmarshalJSON[metaevent.DeviceName](e, t.Event)

		// case "MIDIChannelPrefix":
		//     return unmarshalJSON[metaevent.MIDIChannelPrefix](e, t.Event)

		// case "MIDIPort":
		//     return unmarshalJSON[metaevent.MIDIPort](e, t.Event)

		// case "Tempo":
		//     return unmarshalJSON[metaevent.Tempo](e, t.Event)

		// case "TimeSignature":
		//     return unmarshalJSON[metaevent.TimeSignature](e, t.Event)

		// case "KeySignature":
		//     return unmarshalJSON[metaevent.KeySignature](e, t.Event)

		// case "SMPTEOffset":
		//     return unmarshalJSON[metaevent.SMPTEOffset](e, t.Event)

		// case "EndOfTrack":
		//     return unmarshalJSON[metaevent.EndOfTrack](e, t.Event)

		// case "SequencerSpecificEvent":
		//     return unmarshalJSON[metaevent.SequencerSpecificEvent](e, t.Event)

		// case "NoteOff":
		//     return unmarshalJSON[midievent.NoteOff](e, t.Event)

		// case "NoteOn":
		//     return unmarshalJSON[midievent.NoteOn](e, t.Event)

		// case "PolyphonicPressure":
		//     return unmarshalJSON[midievent.PolyphonicPressure](e, t.Event)

		// case "Controller":
		//     return unmarshalJSON[midievent.Controller](e, t.Event)

		// case "ProgramChange":
		//     return unmarshalJSON[midievent.ProgramChange](e, t.Event)

		// case "ChannelPressure":
		//     return unmarshalJSON[midievent.ChannelPressure](e, t.Event)

		// case "PitchBend":
		//     return unmarshalJSON[midievent.PitchBend](e, t.Event)

		// case "SysExMessage":
		//     return unmarshalJSON[sysex.SysExMessage](e, t.Event)

		// case "SysExContinuation":
		//     return unmarshalJSON[sysex.SysExContinuationMessage](e, t.Event)

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

func equals[T lib.MetaEventType](b byte, t T) bool {
	return (b & 0x7f) == byte(t)
}
