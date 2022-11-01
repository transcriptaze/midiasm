package midifile

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi"
)

type Encoder interface {
	Encode(smf midi.SMF) error
}

type encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) Encoder {
	return &encoder{
		w: w,
	}
}

func (e *encoder) Encode(smf midi.SMF) error {
	if smf.MThd == nil {
		return fmt.Errorf("Missing MThd")
	}

	if int(smf.MThd.Tracks) != len(smf.Tracks) {
		return fmt.Errorf("MThd has incorrect number of tracks")
	}

	if bytes, err := smf.MThd.MarshalBinary(); err != nil {
		return err
	} else if _, err := e.w.Write(bytes); err != nil {
		return err
	}

	return nil
}
