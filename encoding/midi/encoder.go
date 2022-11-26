package midifile

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/lib"
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

	for _, track := range smf.Tracks {
		if bytes, err := EncodeMTrk(*track); err != nil {
			return err
		} else if _, err := e.w.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

func EncodeMTrk(mtrk midi.MTrk) (encoded []byte, err error) {
	type chunk struct {
		delta []byte
		event []byte
	}

	chunks := []chunk{}

	f := func(delta lib.Delta, event encoding.BinaryMarshaler) error {
		if u, err := delta.MarshalBinary(); err != nil {
			return err
		} else if v, err := event.MarshalBinary(); err != nil {
			return err
		} else {
			chunks = append(chunks, chunk{delta: u, event: v})
		}

		return nil
	}

	for _, event := range mtrk.Events {
		delta := lib.Delta(event.Delta())
		if e, ok := event.Event.(encoding.BinaryMarshaler); ok {
			if err = f(delta, e); err != nil {
				return
			}
		} else {
			panic("Expected BinaryMarshaler")
		}
	}

	// ... running status fixup
	var last *chunk
	for i := range chunks {
		chunk := &chunks[i]

		if last != nil {
			status := chunk.event[0]
			if (status&0xf0) == 0x80 || (status&0xf0) == 0x90 {
				if status == last.event[0] {
					chunk.event = chunk.event[1:]
				}
			}
		}

		last = chunk
	}

	// ... get length
	var length uint32
	for _, chunk := range chunks {
		length += uint32(len(chunk.delta))
		length += uint32(len(chunk.event))
	}

	// ... encode
	var b bytes.Buffer

	if _, err = b.Write([]byte(mtrk.Tag)); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, length); err != nil {
		return
	}

	for _, chunk := range chunks {
		if _, err = b.Write(chunk.delta); err != nil {
			return
		}

		if _, err = b.Write(chunk.event); err != nil {
			return
		}
	}

	encoded = b.Bytes()

	return
}
