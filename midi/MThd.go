package midi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MThd struct {
	tag      string
	length   uint32
	format   uint16
	tracks   uint16
	division uint16
	bytes    []byte
}

func (chunk *MThd) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MThd" {
		return fmt.Errorf("Invalid MThd chunk type (%s): expected 'MThd'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])
	if length != 6 {
		return fmt.Errorf("Invalid MThd chunk length (%v): expected 6", length)
	}

	format := binary.BigEndian.Uint16(data[8:10])
	if format != 0 && format != 1 && format != 2 {
		return fmt.Errorf("Invalid MThd chunk format (%v): expected 0,1 or 2", format)
	}

	tracks := binary.BigEndian.Uint16(data[10:12])
	division := binary.BigEndian.Uint16(data[12:14])

	chunk.tag = tag
	chunk.length = length
	chunk.format = format
	chunk.tracks = tracks
	chunk.division = division
	chunk.bytes = data

	return nil
}

func (chunk *MThd) Render(w io.Writer) {
	for _, b := range chunk.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}

	fmt.Fprintf(w, "%12s length:%d format:%d tracks:%d ", chunk.tag, chunk.length, chunk.format, chunk.tracks)
	if chunk.division&0x8000 == 0x0000 {
		fmt.Fprintf(w, "division:metrical time, %d ppqn", chunk.division&0x7fff)
	} else {
		frameSubdivisions := chunk.division & 0x007f
		framesPerSecond := chunk.division & 0x7f00 >> 8

		switch framesPerSecond {
		case 0xe8:
			fmt.Fprintf(w, "division:timecode, %d ticks per SMPE frame, 24 frames per second", frameSubdivisions)
		case 0xe7:
			fmt.Fprintf(w, "division:timecode, %d ticks per SMPE frame, 25 frames per second", frameSubdivisions)
		case 0xe6:
			fmt.Fprintf(w, "division:timecode, %d ticks per SMPE frame, 30 frames per second, drop frame", frameSubdivisions)
		case 0xe5:
			fmt.Fprintf(w, "division:timecode, %d ticks per SMPE frame, 30 frames per second, non-drop frame", frameSubdivisions)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w)
}
