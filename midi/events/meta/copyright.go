package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type Copyright struct {
	event
	Copyright string
}

func NewCopyright(tick uint64, delta uint32, copyright string) (*Copyright, error) {
	return &Copyright{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x02, byte(len(copyright))}, []byte(copyright)...),
			tag:    types.TagCopyright,
			Status: 0xff,
			Type:   0x02,
		},
		Copyright: copyright,
	}, nil
}

func UnmarshalCopyright(tick uint64, delta uint32, bytes []byte) (*Copyright, error) {
	copyright := string(bytes)

	return &Copyright{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  concat([]byte{0x00, 0xff, 0x02, byte(len(bytes))}, bytes),
			tag:    types.TagCopyright,
			Status: 0xff,
			Type:   0x02,
		},
		Copyright: copyright,
	}, nil
}

func (c Copyright) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(c.Status),
		byte(c.Type),
		byte(len(c.Copyright)),
	},
		[]byte(c.Copyright)...), nil
}

func (c *Copyright) UnmarshalText(bytes []byte) error {
	c.tick = 0
	c.delta = 0
	c.bytes = []byte{}
	c.tag = types.TagCopyright
	c.Status = 0xff
	c.Type = 0x02

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Copyright\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Copyright event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		c.delta = uint32(delta)
		c.Copyright = string(match[2])
	}

	return nil
}
