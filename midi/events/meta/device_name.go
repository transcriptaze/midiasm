package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type DeviceName struct {
	event
	Name string
}

func MakeDeviceName(tick uint64, delta uint32, name string) DeviceName {
	return DeviceName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x09, byte(len(name))}, []byte(name)...),
			tag:    types.TagDeviceName,
			Status: 0xff,
			Type:   types.TypeDeviceName,
		},
		Name: name,
	}
}

func UnmarshalDeviceName(tick uint64, delta uint32, bytes []byte) (*DeviceName, error) {
	name := string(bytes)
	event := MakeDeviceName(tick, delta, name)

	return &event, nil
}

func (d DeviceName) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(d.Status),
		byte(d.Type),
		byte(len(d.Name)),
	},
		[]byte(d.Name)...), nil
}

func (d *DeviceName) UnmarshalText(bytes []byte) error {
	d.tick = 0
	d.delta = 0
	d.bytes = []byte{}
	d.tag = types.TagDeviceName
	d.Status = 0xff
	d.Type = types.TypeDeviceName

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)DeviceName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid DeviceName event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		d.delta = uint32(delta)
		d.Name = string(match[2])
	}

	return nil
}
