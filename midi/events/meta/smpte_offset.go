package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/types"
)

type SMPTEOffset struct {
	Tag              string
	Status           types.Status
	Type             types.MetaEventType
	Hour             uint8
	Minute           uint8
	Second           uint8
	FrameRate        uint8
	Frames           uint8
	FractionalFrames uint8
}

func NewSMPTEOffset(bytes []byte) (*SMPTEOffset, error) {
	if len(bytes) != 5 {
		return nil, fmt.Errorf("Invalid SMPTEOffset length (%d): expected '5'", len(bytes))
	}

	rr := (bytes[0] >> 6) & 0x03
	hour := bytes[0] & 0x01f
	minute := bytes[1]
	second := bytes[2]
	frames := bytes[3]
	fractions := bytes[4]

	if hour > 24 {
		return nil, fmt.Errorf("Invalid SMPTEOffset hour (%d): expected a value in the interval [0..24]", hour)
	}

	if minute > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset minute (%d): expected a value in the interval [0..59]", minute)
	}

	if second > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset second (%d): expected a value in the interval [0..59]", second)
	}

	if rr != 0x00 && rr != 0x01 && rr != 0x02 && rr != 0x03 {
		return nil, fmt.Errorf("Invalid SMPTEOffset frame rate (%02X): expected a value in the set [0,1,2,3]", rr)
	}

	framerate := uint8(0)
	switch rr {
	case 0:
		framerate = 24
	case 1:
		framerate = 25
	case 2:
		framerate = 29
	case 3:
		framerate = 30
	}

	if frames >= framerate {
		return nil, fmt.Errorf("Invalid SMPTEOffset frames (%d): expected a value in the interval [0..%d]", frames, framerate-1)
	}

	if fractions > 100 {
		return nil, fmt.Errorf("Invalid SMPTEOffset fractional frames (%d): expected a value in the interval [0..100", fractions)
	}

	return &SMPTEOffset{
		Tag:              "SMPTEOffset",
		Status:           0xff,
		Type:             0x54,
		Hour:             hour,
		Minute:           minute,
		Second:           second,
		FrameRate:        framerate,
		Frames:           frames,
		FractionalFrames: fractions,
	}, nil
}
