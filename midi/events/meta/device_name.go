package metaevent

import (
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type DeviceName struct {
	event
	Name string
}

func MakeDeviceName(tick uint64, delta lib.Delta, name string) DeviceName {
	return DeviceName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x09, byte(len(name))}, []byte(name)...),
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   lib.TypeDeviceName,
		},
		Name: name,
	}
}

func UnmarshalDeviceName(tick uint64, delta lib.Delta, bytes []byte) (*DeviceName, error) {
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
	d.tag = lib.TagDeviceName
	d.Status = 0xff
	d.Type = lib.TypeDeviceName

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)DeviceName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid DeviceName event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		d.delta = delta
		d.Name = string(match[2])
	}

	return nil
}
