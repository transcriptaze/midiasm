package assemble

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/transcriptaze/midiasm/midi"
)

type TextAssembler struct {
}

func NewTextAssembler() TextAssembler {
	return TextAssembler{}
}

func (a TextAssembler) Assemble(source []byte) ([]byte, error) {
	r := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(">> ", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	tracks := make([]*midi.MTrk, 0)

	mthd, err := midi.NewMThd(1, uint16(len(tracks)), 480)
	if err != nil {
		return nil, err
	}

	smf := midi.SMF{
		MThd:   mthd,
		Tracks: tracks,
	}

	return smf.MarshalBinary()
}
