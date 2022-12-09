package events

import (
	"encoding/json"
	"fmt"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
)

func (e *Event) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Event json.RawMessage `json:"event"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	}

	u := struct {
		Tag string `json:"tag"`
	}{}

	if err := json.Unmarshal(t.Event, &u); err != nil {
		return err
	}

	switch u.Tag {
	case "SequenceNumber":
		return unmarshalJSON[metaevent.SequenceNumber](e, t.Event)

	case "Text":
		return unmarshalJSON[metaevent.Text](e, t.Event)

	case "Copyright":
		return unmarshalJSON[metaevent.Copyright](e, t.Event)

	case "TrackName":
		return unmarshalJSON[metaevent.TrackName](e, t.Event)

	case "InstrumentName":
		return unmarshalJSON[metaevent.InstrumentName](e, t.Event)

	case "Lyric":
		return unmarshalJSON[metaevent.Lyric](e, t.Event)

	case "Marker":
		return unmarshalJSON[metaevent.Marker](e, t.Event)

	case "CuePoint":
		return unmarshalJSON[metaevent.CuePoint](e, t.Event)

	case "ProgramName":
		return unmarshalJSON[metaevent.ProgramName](e, t.Event)

	case "DeviceName":
		return unmarshalJSON[metaevent.DeviceName](e, t.Event)

	case "MIDIChannelPrefix":
		return unmarshalJSON[metaevent.MIDIChannelPrefix](e, t.Event)

	case "MIDIPort":
		return unmarshalJSON[metaevent.MIDIPort](e, t.Event)

	case "Tempo":
		return unmarshalJSON[metaevent.Tempo](e, t.Event)

	case "TimeSignature":
		return unmarshalJSON[metaevent.TimeSignature](e, t.Event)

	case "KeySignature":
		return unmarshalJSON[metaevent.KeySignature](e, t.Event)

	case "SMPTEOffset":
		return unmarshalJSON[metaevent.SMPTEOffset](e, t.Event)

	case "EndOfTrack":
		return unmarshalJSON[metaevent.EndOfTrack](e, t.Event)

	case "SequencerSpecificEvent":
		return unmarshalJSON[metaevent.SequencerSpecificEvent](e, t.Event)

	case "NoteOff":
		return unmarshalJSON[midievent.NoteOff](e, t.Event)

	case "NoteOn":
		return unmarshalJSON[midievent.NoteOn](e, t.Event)

	case "PolyphonicPressure":
		return unmarshalJSON[midievent.PolyphonicPressure](e, t.Event)

	case "Controller":
		return unmarshalJSON[midievent.Controller](e, t.Event)

	case "ProgramChange":
		return unmarshalJSON[midievent.ProgramChange](e, t.Event)

	case "ChannelPressure":
		return unmarshalJSON[midievent.ChannelPressure](e, t.Event)

	case "PitchBend":
		return unmarshalJSON[midievent.PitchBend](e, t.Event)

	case "SysExMessage":
		return unmarshalJSON[sysex.SysExMessage](e, t.Event)

	case "SysExContinuation":
		return unmarshalJSON[sysex.SysExContinuationMessage](e, t.Event)

	case "SysExEscape":
		return unmarshalJSON[sysex.SysExEscapeMessage](e, t.Event)

	default:
		return fmt.Errorf("Unrecognised tag (%v)", u.Tag)
	}
}

func unmarshalJSON[
	E TEvent,
	P interface {
		*E
		json.Unmarshaler
	}](e *Event, bytes []byte) (err error) {
	p := P(new(E))

	if err = p.UnmarshalJSON(bytes); err != nil {
		return err
	} else {
		e.Event = *p
	}

	return
}
