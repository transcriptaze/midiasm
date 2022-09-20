package export

import (
	"fmt"
	"io"
	// "github.com/transcriptaze/midiasm/midi/types"
)

type Export struct {
}

func NewExport() (*Export, error) {
	return &Export{}, nil
}

func (x *Export) Export(object interface{}, w io.Writer) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}
