package assemble

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type JSONAssembler struct {
}

type mthd struct {
	Tag    *string `json:"tag,omitempty"`
	Format *uint16 `json:"format,omitempty"`
	PPQN   *uint16 `json:"PPQN,omitempty"`
}

type mtrk struct {
	Tag         *string           `json:"tag,omitempty"`
	TrackNumber *uint16           `json:"tracknumber,omitempty"`
	Events      []json.RawMessage `json:"events"`
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
			mtrk.TrackNumber = lib.TrackNumber(len(smf.Tracks))

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
	for _, bytes := range track.Events {
		e := events.Event{}

		if err := e.UnmarshalJSON(bytes); err != nil {
			return nil, err
		} else {
			mtrk.Events = append(mtrk.Events, &e)
		}
	}

	return fixups(mtrk)
	// return mtrk, fmt.Errorf("missing EndOfTrack")

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
