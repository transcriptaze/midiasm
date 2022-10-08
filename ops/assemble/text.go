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

	mthd := midi.MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   1,
		Tracks:   uint16(len(tracks)),
		Division: 480,
	}

	smf := midi.SMF{
		MThd:   &mthd,
		Tracks: tracks,
	}

	return smf.MarshalBinary()
}
