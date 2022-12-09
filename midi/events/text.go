package events

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
)

func (e *Event) UnmarshalText(text []byte) error {
	s := string(text)

	switch {
	case strings.Contains(s, "SequenceNumber"):
		return unmarshalText[metaevent.SequenceNumber](e, text)

	case strings.Contains(s, "Text"):
		return unmarshalText[metaevent.Text](e, text)

	case strings.Contains(s, "Copyright"):
		return unmarshalText[metaevent.Copyright](e, text)

	case strings.Contains(s, "TrackName"):
		return unmarshalText[metaevent.TrackName](e, text)

	case strings.Contains(s, "InstrumentName"):
		return unmarshalText[metaevent.InstrumentName](e, text)

	case strings.Contains(s, "Lyric"):
		return unmarshalText[metaevent.Lyric](e, text)

	case strings.Contains(s, "Marker"):
		return unmarshalText[metaevent.Marker](e, text)

	case strings.Contains(s, "CuePoint"):
		return unmarshalText[metaevent.CuePoint](e, text)

	case strings.Contains(s, "ProgramName"):
		return unmarshalText[metaevent.ProgramName](e, text)

	case strings.Contains(s, "DeviceName"):
		return unmarshalText[metaevent.DeviceName](e, text)

	case strings.Contains(s, "MIDIChannelPrefix"):
		return unmarshalText[metaevent.MIDIChannelPrefix](e, text)

	case strings.Contains(s, "MIDIPort"):
		return unmarshalText[metaevent.MIDIPort](e, text)

	case strings.Contains(s, "Tempo"):
		return unmarshalText[metaevent.Tempo](e, text)

	case strings.Contains(s, "TimeSignature"):
		return unmarshalText[metaevent.TimeSignature](e, text)

	case strings.Contains(s, "KeySignature"):
		return unmarshalText[metaevent.KeySignature](e, text)

	case strings.Contains(s, "SMPTEOffset"):
		return unmarshalText[metaevent.SMPTEOffset](e, text)

	case strings.Contains(s, "EndOfTrack"):
		return unmarshalText[metaevent.EndOfTrack](e, text)

	case strings.Contains(s, "SequencerSpecificEvent"):
		return unmarshalText[metaevent.SequencerSpecificEvent](e, text)

	case strings.Contains(s, "NoteOff"):
		return unmarshalText[midievent.NoteOff](e, text)

	case strings.Contains(s, "NoteOn"):
		return unmarshalText[midievent.NoteOn](e, text)

	case strings.Contains(s, "PolyphonicPressure"):
		return unmarshalText[midievent.PolyphonicPressure](e, text)

	case strings.Contains(s, "Controller"):
		return unmarshalText[midievent.Controller](e, text)

	case strings.Contains(s, "ProgramChange"):
		return unmarshalText[midievent.ProgramChange](e, text)

	case strings.Contains(s, "ChannelPressure"):
		return unmarshalText[midievent.ChannelPressure](e, text)

	case strings.Contains(s, "PitchBend"):
		return unmarshalText[midievent.PitchBend](e, text)

	case strings.Contains(s, "SysExMessage"):
		return unmarshalText[sysex.SysExMessage](e, text)

	case strings.Contains(s, "SysExContinuation"):
		return unmarshalText[sysex.SysExContinuationMessage](e, text)

	case strings.Contains(s, "SysExEscape"):
		return unmarshalText[sysex.SysExEscapeMessage](e, text)

	default:
		return fmt.Errorf("Unrecognised event (%v)", s)
	}
}

func unmarshalText[
	E TEvent,
	P interface {
		*E
		encoding.TextUnmarshaler
	}](e *Event, text []byte) (err error) {
	p := P(new(E))

	if err = p.UnmarshalText(text); err != nil {
		return err
	} else {
		e.Event = *p
	}

	return
}
