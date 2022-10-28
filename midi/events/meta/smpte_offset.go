package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type SMPTEOffset struct {
	event
	Hour             uint8
	Minute           uint8
	Second           uint8
	FrameRate        uint8
	Frames           uint8
	FractionalFrames uint8
}

func NewSMPTEOffset(tick uint64, delta uint32, hour, minute, second, frameRate, frames, fractionalFrames uint8) (*SMPTEOffset, error) {
	var rr uint8

	if hour > 24 {
		return nil, fmt.Errorf("Invalid SMPTEOffset hour (%d): expected a value in the interval [0..24]", hour)
	}

	if minute > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset minute (%d): expected a value in the interval [0..59]", minute)
	}

	if second > 59 {
		return nil, fmt.Errorf("Invalid SMPTEOffset second (%d): expected a value in the interval [0..59]", second)
	}

	switch frameRate {
	case 24:
		rr = 0 << 6
	case 25:
		rr = 1 << 6
	case 29:
		rr = 2 << 6
	case 30:
		rr = 3 << 6
	default:
		return nil, fmt.Errorf("Invalid SMPTEOffset frame rate (%02X): expected a value in the set [24,25,29,30]", rr)
	}

	if frames >= frameRate {
		return nil, fmt.Errorf("Invalid SMPTEOffset frames (%d): expected a value in the interval [0..%d]", frames, frameRate-1)
	}

	if fractionalFrames > 100 {
		return nil, fmt.Errorf("Invalid SMPTEOffset fractional frames (%d): expected a value in the interval [0..100", fractionalFrames)
	}

	return &SMPTEOffset{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, rr | hour&0x1f, minute, second, frames, fractionalFrames},
			tag:    types.TagSMPTEOffset,
			Status: 0xff,
			Type:   types.TypeSMPTEOffset,
		},
		Hour:             hour,
		Minute:           minute,
		Second:           second,
		FrameRate:        frameRate,
		Frames:           frames,
		FractionalFrames: fractionalFrames,
	}, nil
}

func UnmarshalSMPTEOffset(tick uint64, delta uint32, bytes []byte) (*SMPTEOffset, error) {
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
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  concat([]byte{0x00, 0xff, 0x54, 0x05}, bytes),
			tag:    types.TagSMPTEOffset,
			Status: 0xff,
			Type:   types.TypeSMPTEOffset,
		},
		Hour:             hour,
		Minute:           minute,
		Second:           second,
		FrameRate:        framerate,
		Frames:           frames,
		FractionalFrames: fractions,
	}, nil
}

func (s SMPTEOffset) MarshalBinary() (encoded []byte, err error) {
	var rr uint8
	switch s.FrameRate {
	case 24:
		rr = 0 << 6
	case 25:
		rr = 1 << 6
	case 29:
		rr = 2 << 6
	case 30:
		rr = 3 << 6
	}

	encoded = make([]byte, 8)

	encoded[0] = byte(s.Status)
	encoded[1] = byte(s.Type)
	encoded[2] = byte(5)
	encoded[3] = rr | s.Hour&0x1f
	encoded[4] = s.Minute
	encoded[5] = s.Second
	encoded[6] = s.Frames
	encoded[7] = s.FractionalFrames

	return
}

func (s *SMPTEOffset) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = types.TagSMPTEOffset
	s.Status = 0xff
	s.Type = types.TypeSMPTEOffset

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SMPTEOffset\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 8 {
		return fmt.Errorf("invalid SMPTEOffset event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if hour, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if minute, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if second, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if frameRate, err := strconv.ParseUint(match[5], 10, 8); err != nil {
		return err
	} else if frames, err := strconv.ParseUint(match[6], 10, 8); err != nil {
		return err
	} else if fractions, err := strconv.ParseUint(match[7], 10, 8); err != nil {
		return err
	} else if frameRate != 24 && frameRate != 25 && frameRate != 29 && frameRate != 30 {
		return fmt.Errorf("Invalid SMPTEOffset frame rate (%02X): expected a value in the set [24,25,29,30]", frameRate)
	} else if fractions > 100 {
		return fmt.Errorf("Invalid SMPTEOffset fractional frames (%d): expected a value in the interval [0..100", fractions)
	} else {
		s.delta = uint32(delta)
		s.Hour = uint8(hour)
		s.Minute = uint8(minute)
		s.Second = uint8(second)
		s.FrameRate = uint8(frameRate)
		s.Frames = uint8(frames)
		s.FractionalFrames = uint8(fractions)
	}

	return nil
}
