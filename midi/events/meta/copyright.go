package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type Copyright struct {
	event
	Copyright string
}

func MakeCopyright(tick uint64, delta lib.Delta, copyright string) Copyright {
	return Copyright{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x02, byte(len(copyright))}, []byte(copyright)...),
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
		},
		Copyright: copyright,
	}
}

func UnmarshalCopyright(tick uint64, delta lib.Delta, bytes []byte) (*Copyright, error) {
	copyright := string(bytes)
	event := MakeCopyright(tick, delta, copyright)

	return &event, nil
}

func (c Copyright) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(c.Status),
		byte(c.Type),
		byte(len(c.Copyright)),
	},
		[]byte(c.Copyright)...), nil
}

func (e *Copyright) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Copyright\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Copyright event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagCopyright
		e.Status = 0xff
		e.Type = lib.TypeCopyright
		e.Copyright = string(match[2])
	}

	return nil
}

func (e Copyright) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag       string    `json:"tag"`
		Delta     lib.Delta `json:"delta"`
		Status    byte      `json:"status"`
		Type      byte      `json:"type"`
		Copyright string    `json:"copyright"`
	}{
		Tag:       fmt.Sprintf("%v", e.tag),
		Delta:     e.delta,
		Status:    byte(e.Status),
		Type:      byte(e.Type),
		Copyright: e.Copyright,
	}

	return json.Marshal(t)
}

func (e *Copyright) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag       string    `json:"tag"`
		Delta     lib.Delta `json:"delta"`
		Copyright string    `json:"copyright"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagCopyright) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagCopyright
		e.Type = lib.TypeCopyright
		e.Copyright = t.Copyright
	}

	return nil
}
