package transpose

import (
	"fmt"
	"io"
	// "github.com/transcriptaze/midiasm/midi"
)

type Transpose struct {
	Writer io.Writer
}

func (t *Transpose) Execute(bytes []byte) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}
