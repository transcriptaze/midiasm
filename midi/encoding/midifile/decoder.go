package midifile

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Decoder struct {
}

func NewDecoder() *Decoder {
	decoder := Decoder{}

	return &decoder
}

func (d *Decoder) Decode(b []byte) (*midi.SMF, error) {
	smf := midi.SMF{}
	r := bufio.NewReader(bytes.NewReader(b))

	// Extract header
	chunk, err := readChunk(r)
	if err == nil && chunk != nil {
		var mthd midi.MThd
		if err := mthd.UnmarshalBinary(chunk); err == nil {
			smf.MThd = &mthd
		}
	}

	// Extract tracks
	for err == nil {
		chunk, err = readChunk(r)
		if err == nil && chunk != nil {
			if string(chunk[0:4]) == "MTrk" {
				mtrk := midi.MTrk{
					TrackNumber: types.TrackNumber(len(smf.Tracks)),
				}

				if err := mtrk.UnmarshalBinary(chunk); err != nil {
					return nil, err
				}

				smf.Tracks = append(smf.Tracks, &mtrk)
			}
		}
	}

	if err != io.EOF {
		return nil, err
	}

	if len(smf.Tracks) != int(smf.MThd.Tracks) {
		return nil, fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", smf.MThd.Tracks, len(smf.Tracks))
	}

	return &smf, nil
}

func readChunk(r *bufio.Reader) ([]byte, error) {
	peek, err := r.Peek(8)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(peek[4:8])
	bytes := make([]byte, length+8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
