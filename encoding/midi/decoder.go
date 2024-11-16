package midifile

import (
	"encoding/binary"
	"errors"
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
	chunks := [][]byte{}

	r := chunker{reader: reader}

	for {
		if chunk, err := r.next(); err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		} else if err != nil {
			break
		} else if chunk != nil {
			chunks = append(chunks, chunk)
		}
	}

	// .. extract MThd
	if len(chunks) < 1 {
		return nil, fmt.Errorf("Missing MThd chunk")
	} else {
		var mthd midi.MThd

		if err := mthd.UnmarshalBinary(chunks[0]); err != nil {
			return nil, err
		} else {
			smf.MThd = &mthd
		}
	}

	// .. extract tracks
	if len(chunks) > 1 {
		for _, chunk := range chunks[1:] {
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
	}

	if len(smf.Tracks) != int(smf.MThd.Tracks) {
		return nil, fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", smf.MThd.Tracks, len(smf.Tracks))
	}

	return &smf, nil
}

type chunker struct {
	reader io.Reader
}

func (r *chunker) next() ([]byte, error) {
	buffer := make([]uint8, 8)

	if _, err := io.ReadFull(r.reader, buffer); err != nil {
		return nil, err
	}

	N := binary.BigEndian.Uint32(buffer[4:])
	chunk := make([]byte, N)

	if _, err := io.ReadFull(r.reader, chunk); err != nil {
		return nil, err
	}

	return append(buffer, chunk...), nil
}
