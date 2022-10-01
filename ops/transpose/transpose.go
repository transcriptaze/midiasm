package transpose

import (
	"fmt"
	"io"
	// "github.com/transcriptaze/midiasm/midi"
)

type Transpose struct {
	Writer io.Writer
}

func (t *Transpose) Execute() error {
	return fmt.Errorf("NOT IMPLEMENTED")
}
