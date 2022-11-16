package assemble

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
)

type JSONAssembler struct {
}

type mthd struct {
	Format *uint16 `json:"format,omitempty"`
	PPQN   *uint16 `json:"PPQN,omitempty"`
}

func NewJSONAssembler() JSONAssembler {
	return JSONAssembler{}
}

func (a JSONAssembler) Assemble(r io.Reader) ([]byte, error) {
	src := struct {
		MThd mthd `json:"header"`
	}{}

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&src); err != nil {
		return nil, err
	}

	smf := midi.SMF{}
	if mthd, err := a.parseMThd(src.MThd); err != nil {
		return nil, err
	} else {
		smf.MThd = mthd
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

	if h.Format == nil {
		return nil, fmt.Errorf("missing or invalid 'format' field in MThd")
	} else if *h.Format != 0 && *h.Format != 1 && *h.Format != 2 {
		return nil, fmt.Errorf("invalid 'format' (%v) in MThd", h.Format)
	} else {
		format = *h.Format
	}

	if h.PPQN == nil {
		return nil, fmt.Errorf("missing 'metrical-time' field in MThd")
	} else {
		ppqn = *h.PPQN
	}

	return midi.NewMThd(format, 0, ppqn)
}
