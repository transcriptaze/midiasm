package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/lib"
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

func MakeSMPTEOffset(tick uint64, delta lib.Delta, hour, minute, second, frameRate, frames, fractionalFrames uint8, bytes ...byte) SMPTEOffset {
	var rr uint8

	if hour > 24 {
		panic(fmt.Errorf("Invalid SMPTEOffset hour (%d): expected a value in the interval [0..24]", hour))
	}

	if minute > 59 {
		panic(fmt.Errorf("Invalid SMPTEOffset minute (%d): expected a value in the interval [0..59]", minute))
	}

	if second > 59 {
		panic(fmt.Errorf("Invalid SMPTEOffset second (%d): expected a value in the interval [0..59]", second))
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
		panic(fmt.Errorf("Invalid SMPTEOffset frame rate (%02X): expected a value in the set [24,25,29,30]", rr))
	}

	if frames >= frameRate {
		panic(fmt.Errorf("Invalid SMPTEOffset frames (%d): expected a value in the interval [0..%d]", frames, frameRate-1))
	}

	if fractionalFrames > 100 {
		panic(fmt.Errorf("Invalid SMPTEOffset fractional frames (%d): expected a value in the interval [0..100", fractionalFrames))
	}

	return SMPTEOffset{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
		},
		Hour:             hour,
		Minute:           minute,
		Second:           second,
		FrameRate:        frameRate,
		Frames:           frames,
		FractionalFrames: fractionalFrames,
	}
}

func UnmarshalSMPTEOffset(tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*SMPTEOffset, error) {
	if len(data) != 5 {
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

	event := MakeSMPTEOffset(tick, delta, hour, minute, second, framerate, frames, fractions, bytes...)

	return &event, nil
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

func (e *SMPTEOffset) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.tag = lib.TagSMPTEOffset
	e.Status = 0xff
	e.Type = lib.TypeSMPTEOffset

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SMPTEOffset\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 8 {
		return fmt.Errorf("invalid SMPTEOffset event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
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
		e.delta = delta
		e.Hour = uint8(hour)
		e.Minute = uint8(minute)
		e.Second = uint8(second)
		e.FrameRate = uint8(frameRate)
		e.Frames = uint8(frames)
		e.FractionalFrames = uint8(fractions)
	}

	return nil
}

func (e SMPTEOffset) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag              string    `json:"tag"`
		Delta            lib.Delta `json:"delta"`
		Status           byte      `json:"status"`
		Type             byte      `json:"type"`
		Hour             uint8     `json:"hour"`
		Minute           uint8     `json:"minute"`
		Second           uint8     `json:"second"`
		FrameRate        uint8     `json:"frame-rate"`
		Frames           uint8     `json:"frames"`
		FractionalFrames uint8     `json:"fractional-frames"`
	}{
		Tag:              fmt.Sprintf("%v", e.tag),
		Delta:            e.delta,
		Status:           byte(e.Status),
		Type:             byte(e.Type),
		Hour:             e.Hour,
		Minute:           e.Minute,
		Second:           e.Second,
		FrameRate:        e.FrameRate,
		Frames:           e.Frames,
		FractionalFrames: e.FractionalFrames,
	}

	return json.Marshal(t)
}

func (e *SMPTEOffset) UnmarshalJSON(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = lib.TagSMPTEOffset
	e.Type = lib.TypeSMPTEOffset

	t := struct {
		Tag              string    `json:"tag"`
		Delta            lib.Delta `json:"delta"`
		Hour             uint8     `json:"hour"`
		Minute           uint8     `json:"minute"`
		Second           uint8     `json:"second"`
		FrameRate        uint8     `json:"frame-rate"`
		Frames           uint8     `json:"frames"`
		FractionalFrames uint8     `json:"fractional-frames"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != fmt.Sprintf("%v", lib.TagSMPTEOffset) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.delta = t.Delta
		e.Hour = t.Hour
		e.Minute = t.Minute
		e.Second = t.Second
		e.FrameRate = t.FrameRate
		e.Frames = t.Frames
		e.FractionalFrames = t.FractionalFrames
	}

	return nil
}
