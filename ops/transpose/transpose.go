package transpose

import (
	"bytes"
	"encoding/binary"

	"github.com/transcriptaze/midiasm/midi"
)

type Transpose struct {
}

func (t *Transpose) Execute(smf *midi.SMF, steps int) ([]byte, error) {
	var b bytes.Buffer

	b.Write(smf.MThd.Bytes)

	for _, track := range smf.Tracks {
		track.Transpose(steps)

		b.Write([]byte(track.Tag))
		binary.Write(&b, binary.BigEndian, track.Length)

		for _, event := range track.Events {
			b.Write(event.Bytes())
		}
	}

	return b.Bytes(), nil
}
