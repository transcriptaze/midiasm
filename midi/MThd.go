package midi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type MThd struct {
	Tag           string
	Length        uint32
	Format        uint16
	Tracks        uint16
	Division      uint16
	PPQN          uint16  // TODO make getter/TextUnmarshal
	SMPTETimeCode bool    // TODO make getter/TextUnmarshal
	SubFrames     uint16  // TODO make getter/TextUnmarshal
	FPS           uint8   // TODO make getter/TextUnmarshal
	DropFrame     bool    // TODO make getter/TextUnmarshal
	Bytes         lib.Hex `json:"-"` // TODO make getter/TextUnmarshal
}

func MakeMThd(format uint16, tracks uint16, division uint16, bytes ...byte) MThd {
	if format != 0 && format != 1 && format != 2 {
		panic(fmt.Errorf("Invalid MThd format (%v): expected 0,1 or 2", format))
	}

	if division&0x8000 == 0x8000 {
		fps := division & 0xff00 >> 8
		if fps != 0xe8 && fps != 0xe7 && fps != 0xe3 && fps != 0xe2 {
			panic(fmt.Errorf("Invalid MThd division SMPTE timecode type (%02X): expected 24, 25, 29 or 30 FPS", fps))
		}
	}

	mthd := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   format,
		Tracks:   tracks,
		Division: division,
		Bytes:    bytes,
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
		case 0xe3:
			mthd.FPS = 29
			mthd.DropFrame = true
		case 0xe2:
			mthd.FPS = 30
			mthd.DropFrame = false
		}
	}

	return mthd
}

func (mthd *MThd) UnmarshalBinary(chunk []byte) error {
	if tag := string(chunk[0:4]); tag != "MThd" {
		return fmt.Errorf("invalid MThd chunk - expected:'%v', got:'%v'", "MThd", tag)
	} else {
		mthd.Tag = tag
	}

	if length := binary.BigEndian.Uint32(chunk[4:8]); length != 6 {
		return fmt.Errorf("invalid MThd chunk length - expected:%v, got:%v", 6, length)
	} else {
		mthd.Length = length
	}

	if format := binary.BigEndian.Uint16(chunk[8:10]); format != 0 && format != 1 && format != 2 {
		return fmt.Errorf("invalid MThd format (%v): expected 0,1 or 2", format)
	} else {
		mthd.Format = format
	}

	mthd.Tracks = binary.BigEndian.Uint16(chunk[10:12])
	mthd.Division = binary.BigEndian.Uint16(chunk[12:14])

	if mthd.Division&0x8000 == 0x0000 {
		mthd.PPQN = mthd.Division & 0x7fff
	} else {
		mthd.SMPTETimeCode = true
		switch fps := mthd.Division & 0xff00 >> 8; fps {
		case 0xe8:
			mthd.FPS = 24
		case 0xe7:
			mthd.FPS = 25
		case 0xe3:
			mthd.FPS = 29
			mthd.DropFrame = true
		case 0xe2:
			mthd.FPS = 30
		default:
			return fmt.Errorf("Invalid MThd division SMPTE timecode type (%v): expected 24,25,29 or 30 FPS", fps)
		}

		mthd.SubFrames = mthd.Division & 0x00ff
	}

	mthd.Bytes = slices.Clone(chunk)

	return nil
}

func (mthd MThd) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer

	if _, err = b.Write([]byte(mthd.Tag)); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, mthd.Length); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, mthd.Format); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, mthd.Tracks); err != nil {
		return
	}

	if err = binary.Write(&b, binary.BigEndian, mthd.Division); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}
