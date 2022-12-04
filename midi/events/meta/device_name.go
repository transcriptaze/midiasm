package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type DeviceName struct {
	event
	Name string
}

func MakeDeviceName(tick uint64, delta lib.Delta, name string, bytes ...byte) DeviceName {
	return DeviceName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   lib.TypeDeviceName,
		},
		Name: name,
	}
}

func (e *DeviceName) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	name := string(data)

	*e = MakeDeviceName(tick, delta, name, bytes...)

	return nil
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

func (e DeviceName) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
		Name   string    `json:"name"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Name:   e.Name,
	}

	return json.Marshal(t)
}

func (e *DeviceName) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Name  string    `json:"name"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagDeviceName) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagDeviceName
		e.Type = lib.TypeDeviceName
		e.Name = t.Name
	}

	return nil
}
