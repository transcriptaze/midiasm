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
	type E json.Unmarshaler

	f := func(bytes []byte, e E) error {
		if err := e.UnmarshalJSON(bytes); err != nil {
			return err
		} else {
			mtrk.Events = append(mtrk.Events, events.NewEvent(e))
		}

		return nil
	}

	g := map[string]func() E{
		"SequenceNumber":         func() E { return &metaevent.SequenceNumber{} },
		"Text":                   func() E { return &metaevent.Text{} },
		"Copyright":              func() E { return &metaevent.Copyright{} },
		"TrackName":              func() E { return &metaevent.TrackName{} },
		"InstrumentName":         func() E { return &metaevent.InstrumentName{} },
		"Lyric":                  func() E { return &metaevent.Lyric{} },
		"Marker":                 func() E { return &metaevent.Marker{} },
		"CuePoint":               func() E { return &metaevent.CuePoint{} },
		"ProgramName":            func() E { return &metaevent.ProgramName{} },
		"DeviceName":             func() E { return &metaevent.DeviceName{} },
		"MIDIChannelPrefix":      func() E { return &metaevent.MIDIChannelPrefix{} },
		"MIDIPort":               func() E { return &metaevent.MIDIPort{} },
		"Tempo":                  func() E { return &metaevent.Tempo{} },
		"TimeSignature":          func() E { return &metaevent.TimeSignature{} },
		"KeySignature":           func() E { return &metaevent.KeySignature{} },
		"SMPTEOffset":            func() E { return &metaevent.SMPTEOffset{} },
		"EndOfTrack":             func() E { return &metaevent.EndOfTrack{} },
		"SequencerSpecificEvent": func() E { return &metaevent.SequencerSpecificEvent{} },
		// "ProgramChange":          func() E { return &midievent.ProgramChange{} },
		// "Controller":             func() E { return &midievent.Controller{} },
		// "NoteOn":                 func() E { return &midievent.NoteOn{} },
		// "NoteOff":                func() E { return &midievent.NoteOff{} },
		// "PolyphonicPressure":     func() E { return &midievent.PolyphonicPressure{} },
		// "ChannelPressure":        func() E { return &midievent.ChannelPressure{} },
		// "PitchBend":              func() E { return &midievent.PitchBend{} },
		// "SysExMessage":           func() E { return &sysex.SysExMessage{} },
		// "SysExContinuation":      func() E { return &sysex.SysExContinuationMessage{} },
		// "SysExEscape":            func() E { return &sysex.SysExEscapeMessage{} },
	}

	for _, e := range track.Events {
		t := struct {
			Tag string `json:"tag"`
		}{}

		if err := json.Unmarshal(e.Event, &t); err != nil {
			return nil, err
		}

		if v, ok := g[t.Tag]; ok {
			event := v()
			if err := f(e.Event, event); err != nil {
				return nil, err
			} else if _, ok := event.(*metaevent.EndOfTrack); ok {
				return fixups(mtrk)
			}
		}
	}

	return mtrk, fmt.Errorf("missing EndOfTrack")

}
