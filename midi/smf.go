package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/meta-events"
	"io"
	"os"
)

type SMF struct {
	header *MThd
	tracks []*MTrk
}

func (smf *SMF) UnmarshalBinary(data []byte) error {
	r := bufio.NewReader(bytes.NewReader(data))
	chunks := make([]Chunk, 0)

	chunk := Chunk(nil)
	err := error(nil)

	for err == nil {
		chunk, err = readChunk(r)
		if err == nil && chunk != nil {
			chunks = append(chunks, chunk)
		}
	}

	if err != io.EOF {
		return err
	}

	if len(chunks) == 0 {
		return fmt.Errorf("contains no MIDI chunks")
	}

	header, ok := chunks[0].(*MThd)
	if !ok {
		return fmt.Errorf("invalid MIDI file - expected MThd chunk, got %T", chunks[0])
	}

	if len(chunks[1:]) != int(header.tracks) {
		return fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", header.tracks, len(chunks[1:]))
	}

	tracks := make([]*MTrk, len(chunks[1:]))
	for i, chunk := range chunks[1:] {
		if track, ok := chunk.(*MTrk); ok {
			tracks[i] = track
		} else {
			return fmt.Errorf("invalid MIDI file - expected MTrk chunk, got %T", chunk)
		}
	}

	smf.header = header
	smf.tracks = tracks

	return nil
}

func (smf *SMF) Render() {
	smf.header.Render(os.Stdout)
	for _, track := range smf.tracks {
		track.Render(os.Stdout)
	}
}

func (smf *SMF) Notes() error {
	tempo := make([]*metaevent.Tempo, 0)
	for _, e := range smf.tracks[0].Events {
		if v, ok := e.(*metaevent.Tempo); ok {
			tempo = append(tempo, v)
		}
	}

	for _, track := range smf.tracks[1:] {
		track.Notes(smf.header.division, tempo, os.Stdout)
	}

	return nil
}

func readChunk(r *bufio.Reader) (Chunk, error) {
	peek, err := r.Peek(8)
	if err != nil {
		return nil, err
	}

	tag := string(peek[0:4])
	length := binary.BigEndian.Uint32(peek[4:8])

	bytes := make([]byte, length+8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}

	switch tag {
	case "MThd":
		var mthd MThd
		if err := mthd.UnmarshalBinary(bytes); err != nil {
			return nil, err
		}
		return &mthd, nil

	case "MTrk":
		var mtrk MTrk
		if err := mtrk.UnmarshalBinary(bytes); err != nil {
			return nil, err
		}
		return &mtrk, nil
	}

	return nil, nil
}
