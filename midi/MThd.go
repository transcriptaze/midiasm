package midi

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/transcriptaze/midiasm/midi/types"
)

type MThd struct {
	Tag           string
	Length        uint32
	Format        uint16
	Tracks        uint16
	Division      uint16
	PPQN          uint16    // TODO make getter/TextUnmarshal
	SMPTETimeCode bool      // TODO make getter/TextUnmarshal
	SubFrames     uint16    // TODO make getter/TextUnmarshal
	FPS           uint8     // TODO make getter/TextUnmarshal
	DropFrame     bool      // TODO make getter/TextUnmarshal
	Bytes         types.Hex `json:"-"` // TODO make getter/TextUnmarshal
}

func NewMThd(format uint16, tracks uint16, division uint16) (*MThd, error) {
	if format != 0 && format != 1 && format != 2 {
		return nil, fmt.Errorf("Invalid MThd format (%v): expected 0,1 or 2", format)
	}

	if division&0x8000 == 0x8000 {
		fps := division & 0xff00 >> 8
		if fps != 0xe8 && fps != 0xe7 && fps != 0xe6 && fps != 0xe5 {
			return nil, fmt.Errorf("Invalid MThd division SMPTE timecode type (%02X): expected E8, E7, E6 or E5", fps)
		}
	}

	mthd := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   format,
		Tracks:   tracks,
		Division: division,
	}

	if bytes, err := mthd.MarshalBinary(); err != nil {
		return nil, err
	} else {
		mthd.Bytes = bytes
	}

	if division&0x8000 == 0x0000 {
		mthd.SMPTETimeCode = false
		mthd.PPQN = division & 0x7fff
	} else {
		mthd.SMPTETimeCode = true
		mthd.SubFrames = division & 0x007f

		fps := division & 0xff00 >> 8
		switch fps {
		case 0xe8:
			mthd.FPS = 24
			mthd.DropFrame = false
		case 0xe7:
			mthd.FPS = 25
			mthd.DropFrame = false
		case 0xe6:
			mthd.FPS = 30
			mthd.DropFrame = true
		case 0xe5:
			mthd.FPS = 30
			mthd.DropFrame = false
		}
	}

	return &mthd, nil
}

func (chunk MThd) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer

	if _, err = b.Write([]byte(chunk.Tag)); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, chunk.Length); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, chunk.Format); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, chunk.Tracks); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, chunk.Division); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}

func (chunk *MThd) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MThd" {
		return fmt.Errorf("Invalid MThd chunk type (%s): expected 'MThd'", tag)
	}

	if len(data) < 14 {
		return fmt.Errorf("Insufficent bytes in  MThd (%v): expected 14", len(data))
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

	if division&0x8000 == 0x8000 {
		fps := division & 0xff00 >> 8
		if fps != 0xe8 && fps != 0xe7 && fps != 0xe6 && fps != 0xe5 {
			return fmt.Errorf("Invalid MThd division SMPTE timecode type (%02X): expected E8, E7, E6 or E5", fps)
		}
	}

	chunk.Tag = tag
	chunk.Length = length
	chunk.Format = format
	chunk.Tracks = tracks
	chunk.Division = division
	chunk.Bytes = data

	if division&0x8000 == 0x0000 {
		chunk.SMPTETimeCode = false
		chunk.PPQN = chunk.Division & 0x7fff
	} else {
		chunk.SMPTETimeCode = true
		chunk.SubFrames = division & 0x007f

		fps := division & 0xff00 >> 8
		switch fps {
		case 0xe8:
			chunk.FPS = 24
			chunk.DropFrame = false
		case 0xe7:
			chunk.FPS = 25
			chunk.DropFrame = false
		case 0xe6:
			chunk.FPS = 30
			chunk.DropFrame = true
		case 0xe5:
			chunk.FPS = 30
			chunk.DropFrame = false
		}
	}
	return nil
}
