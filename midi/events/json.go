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
		return unmarshal[metaevent.SequenceNumber](e, t.Event)

	case "Text":
		return unmarshal[metaevent.Text](e, t.Event)

	case "Copyright":
		return unmarshal[metaevent.Copyright](e, t.Event)

	case "TrackName":
		return unmarshal[metaevent.TrackName](e, t.Event)

	case "InstrumentName":
		return unmarshal[metaevent.InstrumentName](e, t.Event)

	case "Lyric":
		return unmarshal[metaevent.Lyric](e, t.Event)

	case "Marker":
		return unmarshal[metaevent.Marker](e, t.Event)

	case "CuePoint":
		return unmarshal[metaevent.CuePoint](e, t.Event)

	case "ProgramName":
		return unmarshal[metaevent.ProgramName](e, t.Event)

	case "DeviceName":
		return unmarshal[metaevent.DeviceName](e, t.Event)

	case "MIDIChannelPrefix":
		return unmarshal[metaevent.MIDIChannelPrefix](e, t.Event)

	case "MIDIPort":
		return unmarshal[metaevent.MIDIPort](e, t.Event)

	case "Tempo":
		return unmarshal[metaevent.Tempo](e, t.Event)

	case "TimeSignature":
		return unmarshal[metaevent.TimeSignature](e, t.Event)

	case "KeySignature":
		return unmarshal[metaevent.KeySignature](e, t.Event)

	case "SMPTEOffset":
		return unmarshal[metaevent.SMPTEOffset](e, t.Event)

	case "EndOfTrack":
		return unmarshal[metaevent.EndOfTrack](e, t.Event)

	case "SequencerSpecificEvent":
		return unmarshal[metaevent.SequencerSpecificEvent](e, t.Event)

	case "NoteOff":
		return unmarshal[midievent.NoteOff](e, t.Event)

	case "NoteOn":
		return unmarshal[midievent.NoteOn](e, t.Event)

	case "PolyphonicPressure":
		return unmarshal[midievent.PolyphonicPressure](e, t.Event)

	case "Controller":
		return unmarshal[midievent.Controller](e, t.Event)

	case "ProgramChange":
		return unmarshal[midievent.ProgramChange](e, t.Event)

	case "ChannelPressure":
		return unmarshal[midievent.ChannelPressure](e, t.Event)

	case "PitchBend":
		return unmarshal[midievent.PitchBend](e, t.Event)

	case "SysExMessage":
		return unmarshal[sysex.SysExMessage](e, t.Event)

	case "SysExContinuation":
		return unmarshal[sysex.SysExContinuationMessage](e, t.Event)

	case "SysExEscape":
		return unmarshal[sysex.SysExEscapeMessage](e, t.Event)

	default:
		return fmt.Errorf("Unrecognised tag (%v)", u.Tag)
	}
}

func unmarshal[
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
