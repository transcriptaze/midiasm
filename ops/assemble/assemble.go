package assemble

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi"
)

type Assembler interface {
	Assemble([]byte) ([]byte, error)
}

func assemble(smf midi.SMF) ([]byte, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
