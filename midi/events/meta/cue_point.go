package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type CuePoint struct {
	event
	CuePoint string
}

func MakeCuePoint(tick uint64, delta uint32, cuepoint string) CuePoint {
	return CuePoint{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x07, byte(len(cuepoint))}, []byte(cuepoint)...),
			tag:    types.TagCuePoint,
			Status: 0xff,
			Type:   types.TypeCuePoint,
		},
		CuePoint: cuepoint,
	}
}

func UnmarshalCuePoint(tick uint64, delta uint32, bytes []byte) (*CuePoint, error) {
	cuepoint := string(bytes)
	event := MakeCuePoint(tick, delta, cuepoint)

	return &event, nil
}

func (c CuePoint) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(c.Status),
		byte(c.Type),
		byte(len(c.CuePoint)),
	},
		[]byte(c.CuePoint)...), nil
}

func (c *CuePoint) UnmarshalText(bytes []byte) error {
	c.tick = 0
	c.delta = 0
	c.bytes = []byte{}
	c.tag = types.TagCuePoint
	c.Status = 0xff
	c.Type = types.TypeCuePoint

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)CuePoint\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid CuePoint event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		c.delta = uint32(delta)
		c.CuePoint = string(match[2])
	}

	return nil
}
