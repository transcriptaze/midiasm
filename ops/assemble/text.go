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

	mthd := midi.MThd{
		Tag:    "MThd",
		Length: 6,
		// Format        : ,
		Tracks: 2,
		// Division      : ,
		// PPQN          : ,
		// SMPTETimeCode : ,
		// SubFrames     : ,
		// FPS           : ,
		// DropFrame     : ,
	}

	tracks := make([]*midi.MTrk, mthd.Tracks)

	smf := midi.SMF{
		MThd:   &mthd,
		Tracks: tracks,
	}

	return assemble(smf)
}
