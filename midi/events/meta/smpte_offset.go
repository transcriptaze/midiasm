package metaevent

import (
	"fmt"
	"io"
)

type SMPTEOffset struct {
	MetaEvent
	FrameRate        uint8
	Hour             uint8
	Minute           uint8
	Second           uint8
	Frames           uint8
	FractionalFrames uint8
}

func NewSMPTEOffset(event *MetaEvent, r io.ByteReader) (*SMPTEOffset, error) {
	if event.Type != 0x58 {
		return nil, fmt.Errorf("Invalid SMPTEOffset event type (%02x): expected '58'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 5 {
		return nil, fmt.Errorf("Invalid SMPTEOffset length (%d): expected '5'", len(data))
	}

	framerate := (data[0] >> 5) & 0x03
	hour := data[0] & 0x01f
	minute := data[1]
	second := data[2]
	frames := data[3]
	fractions := data[4]

	if framerate != 0x00 && framerate != 0x01 && framerate != 0x10 && framerate != 0x11 {
		return nil, fmt.Errorf("Invalid SMPTEOffset frame rate (%02X): expected a value in the set [00,01,10,11]", hour)
	}

	if hour > 24 {
		return nil, fmt.Errorf("Invalid SMPTEOffset hour (%d): expected a value in the interval [0..24]", hour)
	}

	if minute > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset minute (%d): expected a value in the interval [0..59]", minute)
	}

	if second > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset second (%d): expected a value in the interval [0..59]", second)
	}

	if frames > 23 && frames != 24 && frames != 28 && frames != 29 {
		return nil, fmt.Errorf("Invalid SMPTEOffset frames (%d): expected a value in the interval [0..23,24,28,29", second)
	}

	return &SMPTEOffset{
		MetaEvent:        *event,
		FrameRate:        framerate,
		Hour:             hour,
		Second:           minute,
		Minute:           second,
		Frames:           frames,
		FractionalFrames: fractions,
	}, nil
}

func (e *SMPTEOffset) Render(w io.Writer) {
	framerate := "24 fps"
	switch e.FrameRate {
	case 0x00:
		framerate = "24fps"
	case 0x01:
		framerate = "25fps"
	case 0x10:
		framerate = "30fps (drop frame)"
	case 0x11:
		framerate = "30fps (non-drop frame)"
	}

	fmt.Fprintf(w, "%s %-16s %s, %02d:%02d:%02d, %d frames, %d fractional frames", e.MetaEvent, "SMPTEOffset", framerate, e.Hour, e.Minute, e.Second, e.Frames, e.FractionalFrames)
}
