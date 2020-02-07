package metaevent

import (
	"fmt"
	"io"
)

type SMPTEOffset struct {
	Tag string
	MetaEvent
	Hour             uint8
	Minute           uint8
	Second           uint8
	FrameRate        uint8
	Frames           uint8
	FractionalFrames uint8
}

func NewSMPTEOffset(event *MetaEvent, r io.ByteReader) (*SMPTEOffset, error) {
	if event.Type != 0x54 {
		return nil, fmt.Errorf("Invalid SMPTEOffset event type (%02x): expected '54'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 5 {
		return nil, fmt.Errorf("Invalid SMPTEOffset length (%d): expected '5'", len(data))
	}

	rr := (data[0] >> 6) & 0x03
	hour := data[0] & 0x01f
	minute := data[1]
	second := data[2]
	frames := data[3]
	fractions := data[4]

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
		MetaEvent:        *event,
		Hour:             hour,
		Minute:           minute,
		Second:           second,
		FrameRate:        framerate,
		Frames:           frames,
		FractionalFrames: fractions,
	}, nil
}
