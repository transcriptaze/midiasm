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
		fmt.Fprintf(w, "%02X ", b)
	}

	fmt.Fprintf(w, "  %s length:%d, format:%d, tracks:%d ", chunk.tag, chunk.length, chunk.format, chunk.tracks)
	if chunk.division&0x8000 == 0x0000 {
		fmt.Fprintf(w, ", metrical time, %d ppqn", chunk.division&0x7fff)
	} else {
		fps := chunk.division & 0xff00 >> 8
		subdivisions := chunk.division & 0x007f

		switch fps {
		case 0xe8:
			fmt.Fprintf(w, ", SMPTE timecode, %d ticks per frame, 24 fps", subdivisions)
		case 0xe7:
			fmt.Fprintf(w, ", SMPTE timecode, %d ticks per frame, 25 fps", subdivisions)
		case 0xe6:
			fmt.Fprintf(w, ", SMPTE timecode, %d ticks per frame, 30 fps (drop frame)", subdivisions)
		case 0xe5:
			fmt.Fprintf(w, ", SMPTE timecode, %d ticks per frame, 30 fps (non-drop frame)", subdivisions)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w)
}
