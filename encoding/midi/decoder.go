package midifile

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type Decoder interface {
	Decode(reader io.Reader) (*midi.SMF, error)
}

type decoder struct {
}

func NewDecoder() Decoder {
	return &decoder{}
}

func (d *decoder) Decode(reader io.Reader) (*midi.SMF, error) {
	smf := midi.SMF{}
	chunks := make(chan []byte)
	errors := make(chan error)

	defer close(errors)

	go d.read(reader, chunks, errors)

	list := [][]byte{}

loop:
	for {
		select {
		case chunk, ok := <-chunks:
			if ok {
				list = append(list, chunk)
			} else {
				break loop
			}

		case err := <-errors:
			return nil, err
		}
	}

	// .. extract header

	for _, chunk := range list {
		var mthd midi.MThd
		if err := mthd.UnmarshalBinary(chunk); err == nil {
			smf.MThd = &mthd
		}
		break
	}

	// .. extract tracks
	for _, chunk := range list[1:] {
		if string(chunk[0:4]) == "MTrk" {
			mtrk := midi.MTrk{
				TrackNumber: lib.TrackNumber(len(smf.Tracks)),
				Context:     context.NewContext(),
			}

			if err := mtrk.UnmarshalBinary(chunk); err != nil {
				return nil, err
			}

			smf.Tracks = append(smf.Tracks, &mtrk)
		}
	}

	if len(smf.Tracks) != int(smf.MThd.Tracks) {
		return nil, fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", smf.MThd.Tracks, len(smf.Tracks))
	}

	return &smf, nil
}

func (d *decoder) read(reader io.Reader, chunks chan []byte, errors chan error) {
	defer close(chunks)

	r := bufio.NewReader(reader)

	for {
		peek, err := r.Peek(8)
		switch {
		case err != nil && err != io.EOF:
			errors <- err
			return

		case err != nil && err == io.EOF:
			return

		default:
			length := binary.BigEndian.Uint32(peek[4:8])
			chunk := make([]byte, length+8)
			if _, err := io.ReadFull(r, chunk); err != nil {
				errors <- err
				return
			} else {
				chunks <- chunk
			}
		}
	}
}
