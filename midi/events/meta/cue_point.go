package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type CuePoint struct {
	event
	CuePoint string
}

func MakeCuePoint(tick uint64, delta lib.Delta, cuepoint string, bytes ...byte) CuePoint {
	return CuePoint{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagCuePoint,
			Status: 0xff,
			Type:   lib.TypeCuePoint,
		},
		CuePoint: cuepoint,
	}
}

func UnmarshalCuePoint(ctx *context.Context, tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*CuePoint, error) {
	cuepoint := string(data)
	event := MakeCuePoint(tick, delta, cuepoint, bytes...)

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
	c.tag = lib.TagCuePoint
	c.Status = 0xff
	c.Type = lib.TypeCuePoint

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)CuePoint\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid CuePoint event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		c.delta = delta
		c.CuePoint = string(match[2])
	}

	return nil
}

func (e CuePoint) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag      string    `json:"tag"`
		Delta    lib.Delta `json:"delta"`
		Status   byte      `json:"status"`
		Type     byte      `json:"type"`
		CuePoint string    `json:"cuepoint"`
	}{
		Tag:      fmt.Sprintf("%v", e.tag),
		Delta:    e.delta,
		Status:   byte(e.Status),
		Type:     byte(e.Type),
		CuePoint: e.CuePoint,
	}

	return json.Marshal(t)
}

func (e *CuePoint) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string    `json:"tag"`
		Delta    lib.Delta `json:"delta"`
		CuePoint string    `json:"cuepoint"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagCuePoint) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagCuePoint
		e.Type = lib.TypeCuePoint
		e.CuePoint = t.CuePoint
	}

	return nil
}
