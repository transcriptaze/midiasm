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
	ix := 0

	for ix < len(chunks) {
		chunk := chunks[ix]
		ix++

		tag := string(chunk[0:4])
		length := binary.BigEndian.Uint32(chunk[4:8])

		if tag == "MThd" && length >= 6 {
			format := binary.BigEndian.Uint16(chunk[8:10])
			tracks := binary.BigEndian.Uint16(chunk[10:12])
			division := binary.BigEndian.Uint16(chunk[12:14])

			if format != 0 && format != 1 && format != 2 {
				return nil, fmt.Errorf("Invalid MThd format (%v): expected 0,1 or 2", format)
			}

			if division&0x8000 == 0x8000 {
				fps := division & 0xff00 >> 8
				if fps != 0xe8 && fps != 0xe7 && fps != 0xe3 && fps != 0xe2 {
					return nil, fmt.Errorf("Invalid MThd division SMPTE timecode type (%02X): expected E8, E7, E3 or E2", fps)
				}
			}

			mthd := midi.MakeMThd(format, tracks, division, chunk...)
			smf.MThd = &mthd
			break
		}
	}

	// .. extract tracks
	for ix < len(chunks) {
		chunk := chunks[ix]
		ix++

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
