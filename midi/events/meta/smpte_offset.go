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
	if hour > 24 {
		panic(fmt.Errorf("Invalid SMPTE offset hour (%d): expected a value in the interval [0..24]", hour))
	}

	if minute > 59 {
		panic(fmt.Errorf("Invalid SMPTE offset minute (%d): expected a value in the interval [0..59]", minute))
	}

	if second > 59 {
		panic(fmt.Errorf("Invalid SMPTE offset second (%d): expected a value in the interval [0..59]", second))
	}

	if frameRate != 24 && frameRate != 25 && frameRate != 29 && frameRate != 30 {
		panic(fmt.Errorf("Invalid SMPTE offset frame rate (%02X): expected a value in the set [24,25,29,30]", frameRate))
	}

	if frames >= frameRate {
		panic(fmt.Errorf("Invalid SMPTE offset frames (%d): expected a value in the interval [0..%d]", frames, frameRate-1))
	}

	if fractionalFrames > 100 {
		panic(fmt.Errorf("Invalid SMPTE offset fractional frames (%d): expected a value in the interval [0..100", fractionalFrames))
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

func (e *SMPTEOffset) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 5 {
		return fmt.Errorf("Invalid SMPTE offset length (%d): expected '5'", len(data))
	}

	rr := (data[0] >> 5) & 0x03
	hour := data[0] & 0x01f
	minute := data[1]
	second := data[2]
	frames := data[3]
	fractions := data[4]

	if hour > 24 {
		return fmt.Errorf("Invalid SMPTE offset hour (%d): expected a value in the interval [0..24]", hour)
	}

	if minute > 59 {
		return fmt.Errorf("Invalid SMPTE offset minute (%d): expected a value in the interval [0..59]", minute)
	}

	if second > 59 {
		return fmt.Errorf("Invalid SMPTE offset second (%d): expected a value in the interval [0..59]", second)
	}

	if rr != 0x00 && rr != 0x01 && rr != 0x02 && rr != 0x03 {
		return fmt.Errorf("Invalid SMPTE offset frame rate (%02X): expected a value in the set [0,1,2,3]", rr)
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
		return fmt.Errorf("Invalid SMPTE offset frames (%d): expected a value in the interval [0..%d]", frames, framerate-1)
	}

	if fractions > 100 {
		return fmt.Errorf("Invalid SMPTE offset fractional frames (%d): expected a value in the interval [0..100", fractions)
	}

	*e = MakeSMPTEOffset(tick, delta, hour, minute, second, framerate, frames, fractions, bytes...)

	return nil
}

func (s SMPTEOffset) MarshalBinary() (encoded []byte, err error) {
	var rr uint8
	switch s.FrameRate {
	case 24:
		rr = 0 << 5
	case 25:
		rr = 1 << 5
	case 29:
		rr = 2 << 5
	case 30:
		rr = 3 << 5
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

func (e *SMPTEOffset) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagSMPTEOffset, remaining[0])
	} else if !equals(remaining[1], lib.TypeSMPTEOffset) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagSMPTEOffset, remaining[1])
	} else if data, err := vlf(remaining[2:]); err != nil {
		return err
	} else if len(data) < 5 {
		return fmt.Errorf("Invalid SMPTE offset data")
	} else {
		rr := (data[0] >> 5) & 0x03
		hour := data[0] & 0x01f
		minute := data[1]
		second := data[2]
		frames := data[3]
		fractions := data[4]

		if hour > 24 {
			return fmt.Errorf("Invalid SMPTE offset hour (%d): expected a value in the interval [0..24]", hour)
		}

		if minute > 59 {
			return fmt.Errorf("Invalid SMPTE offset minute (%d): expected a value in the interval [0..59]", minute)
		}

		if second > 59 {
			return fmt.Errorf("Invalid SMPTE offset second (%d): expected a value in the interval [0..59]", second)
		}

		if rr != 0x00 && rr != 0x01 && rr != 0x02 && rr != 0x03 {
			return fmt.Errorf("Invalid SMPTE offset frame rate (%02X): expected a value in the set [0,1,2,3]", rr)
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
			return fmt.Errorf("Invalid SMPTE offset frames (%d): expected a value in the interval [0..%d]", frames, framerate-1)
		}

		if fractions > 100 {
			return fmt.Errorf("Invalid SMPTE offset fractional frames (%d): expected a value in the interval [0..100", fractions)
		}

		*e = MakeSMPTEOffset(0, delta, hour, minute, second, framerate, frames, fractions, bytes...)
	}

	return nil
}

func (e *SMPTEOffset) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SMPTEOffset\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 8 {
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
		return fmt.Errorf("Invalid SMPTE offset frame rate (%02X): expected a value in the set [24,25,29,30]", frameRate)
	} else if fractions > 100 {
		return fmt.Errorf("Invalid SMPTE offset fractional frames (%d): expected a value in the interval [0..100", fractions)
	} else {
		*e = MakeSMPTEOffset(0, delta, uint8(hour), uint8(minute), uint8(second), uint8(frameRate), uint8(frames), uint8(fractions), []byte{}...)
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
		*e = MakeSMPTEOffset(0, t.Delta, t.Hour, t.Minute, t.Second, t.FrameRate, t.Frames, t.FractionalFrames, []byte{}...)
	}

	return nil
}
