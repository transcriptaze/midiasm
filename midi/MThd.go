package midi

import (
	"bytes"
	"encoding/binary"
	"fmt"

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
		if fps != 0xe8 && fps != 0xe7 && fps != 0xe6 && fps != 0xe5 {
			panic(fmt.Errorf("Invalid MThd division SMPTE timecode type (%02X): expected E8, E7, E6 or E5", fps))
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
		case 0xe6:
			mthd.FPS = 30
			mthd.DropFrame = true
		case 0xe5:
			mthd.FPS = 30
			mthd.DropFrame = false
		}
	}

	return mthd
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
