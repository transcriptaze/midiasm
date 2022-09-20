package export

import (
	"encoding/json"
	"io"

	"github.com/transcriptaze/midiasm/midi"
)

type Export struct {
}

func NewExport() (*Export, error) {
	return &Export{}, nil
}

func (x *Export) Export(smf *midi.SMF, w io.Writer) error {
	if bytes, err := json.MarshalIndent(smf, "", "  "); err != nil {
		return err
	} else if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}
