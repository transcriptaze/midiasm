package assemble

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
)

type JSONAssembler struct {
}

type mthd struct {
	Tag    *string `json:"tag,omitempty"`
	Format *uint16 `json:"format,omitempty"`
	PPQN   *uint16 `json:"PPQN,omitempty"`
}

type mtrk struct {
	Tag         *string `json:"tag,omitempty"`
	TrackNumber *uint16 `json:"tracknumber,omitempty"`
	Events      []struct {
		Event json.RawMessage `json:"event"`
	} `json:"events"`
}

func NewJSONAssembler() JSONAssembler {
	return JSONAssembler{}
}

func (a JSONAssembler) Assemble(r io.Reader) ([]byte, error) {
	src := struct {
		Header mthd   `json:"header"`
		Tracks []mtrk `json:"tracks"`
	}{}

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&src); err != nil {
		return nil, err
	}

	smf := midi.SMF{}

	// ... header
	if mthd, err := a.parseMThd(src.Header); err != nil {
		return nil, err
	} else {
		smf.MThd = mthd
	}

	// ... tracks

	for _, t := range src.Tracks {
		if mtrk, err := a.parseMTrk(t); err != nil {
			return nil, err
		} else {
			smf.Tracks = append(smf.Tracks, mtrk)
			smf.MThd.Tracks += 1
		}
	}

	// ... assemble into MIDI file
	var b bytes.Buffer
	var e = midifile.NewEncoder(&b)

	if err := e.Encode(smf); err != nil {
		return nil, err
	} else {
		return b.Bytes(), nil
	}
}

func (a JSONAssembler) parseMThd(h mthd) (*midi.MThd, error) {
	var format uint16
	var ppqn uint16

	if h.Tag == nil || *h.Tag != "MThd" {
		return nil, fmt.Errorf("missing or invalid 'MThd' tag in header")
	} else if h.Format == nil {
		return nil, fmt.Errorf("missing or invalid 'format' field in header")
	} else if *h.Format != 0 && *h.Format != 1 && *h.Format != 2 {
		return nil, fmt.Errorf("invalid 'format' (%v) in header", h.Format)
	} else {
		format = *h.Format
	}

	if h.PPQN == nil {
		return nil, fmt.Errorf("missing 'metrical-time' field in header")
	} else {
		ppqn = *h.PPQN
	}

	return midi.NewMThd(format, 0, ppqn)
}

func (a JSONAssembler) parseMTrk(track mtrk) (*midi.MTrk, error) {
	var mtrk *midi.MTrk

	// ... MTrk header
	if track.Tag == nil || *track.Tag != "MTrk" {
		return nil, fmt.Errorf("missing or invalid 'MTrk' tag in track")
	} else if v, err := midi.NewMTrk(); err != nil {
		return nil, err
	} else if v == nil {
		return nil, fmt.Errorf("error creating 'MTrk' for tracks")
	} else {
		mtrk = v
	}

	// ... events
	for _, e := range track.Events {
		t := struct {
			Tag string `json:"tag"`
		}{}

		if err := json.Unmarshal(e.Event, &t); err != nil {
			return nil, err
		}

		switch t.Tag {
		case "SequenceNumber":
			if err := unmarshal[metaevent.SequenceNumber](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Text":
			if err := unmarshal[metaevent.Text](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Copyright":
			if err := unmarshal[metaevent.Copyright](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "TrackName":
			if err := unmarshal[metaevent.TrackName](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "InstrumentName":
			if err := unmarshal[metaevent.InstrumentName](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Lyric":
			if err := unmarshal[metaevent.Lyric](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Marker":
			if err := unmarshal[metaevent.Marker](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "CuePoint":
			if err := unmarshal[metaevent.CuePoint](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "ProgramName":
			if err := unmarshal[metaevent.ProgramName](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "DeviceName":
			if err := unmarshal[metaevent.DeviceName](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "MIDIChannelPrefix":
			if err := unmarshal[metaevent.MIDIChannelPrefix](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "MIDIPort":
			if err := unmarshal[metaevent.MIDIPort](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Tempo":
			if err := unmarshal[metaevent.Tempo](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "TimeSignature":
			if err := unmarshal[metaevent.TimeSignature](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "KeySignature":
			if err := unmarshal[metaevent.KeySignature](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "SMPTEOffset":
			if err := unmarshal[metaevent.SMPTEOffset](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "EndOfTrack":
			if err := unmarshal[metaevent.EndOfTrack](mtrk, e.Event); err != nil {
				return nil, err
			} else {
				return fixups(mtrk)
			}

		case "SequencerSpecificEvent":
			if err := unmarshal[metaevent.SequencerSpecificEvent](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "NoteOff":
			if err := unmarshal[midievent.NoteOff](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "NoteOn":
			if err := unmarshal[midievent.NoteOn](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "PolyphonicPressure":
			if err := unmarshal[midievent.PolyphonicPressure](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "Controller":
			if err := unmarshal[midievent.Controller](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "ProgramChange":
			if err := unmarshal[midievent.ProgramChange](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "ChannelPressure":
			if err := unmarshal[midievent.ChannelPressure](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "PitchBend":
			if err := unmarshal[midievent.PitchBend](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "SysExMessage":
			if err := unmarshal[sysex.SysExMessage](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "SysExContinuation":
			if err := unmarshal[sysex.SysExContinuationMessage](mtrk, e.Event); err != nil {
				return nil, err
			}

		case "SysExEscape":
			if err := unmarshal[sysex.SysExEscapeMessage](mtrk, e.Event); err != nil {
				return nil, err
			}
		}
	}

	return mtrk, fmt.Errorf("missing EndOfTrack")

}

func unmarshal[
	E events.TEvent,
	P interface {
		*E
		json.Unmarshaler
	}](mtrk *midi.MTrk, bytes []byte) (err error) {
	p := P(new(E))

	if err = p.UnmarshalJSON(bytes); err == nil {
		mtrk.Events = append(mtrk.Events, events.NewEvent(*p))
	}

	return
}
