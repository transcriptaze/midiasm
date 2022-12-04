package metaevent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TrackName struct {
	event
	Name string
}

func MakeTrackName(tick uint64, delta lib.Delta, name string, bytes ...byte) TrackName {
	return TrackName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
		},
		Name: name,
	}
}

func UnmarshalTrackName(ctx *context.Context, tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*TrackName, error) {
	name := string(data)
	event := MakeTrackName(tick, delta, name, bytes...)

	return &event, nil
}

func (t TrackName) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer
	var v []byte

	if err = b.WriteByte(byte(t.Status)); err != nil {
		return
	}

	if err = b.WriteByte(byte(t.Type)); err != nil {
		return
	}

	name := lib.VLF(t.Name)
	if v, err = name.MarshalBinary(); err != nil {
		return
	} else if _, err = b.Write(v); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}

func (e *TrackName) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.tag = lib.TagTrackName
	e.Status = 0xff
	e.Type = lib.TypeTrackName

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)TrackName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid %v event (%v)", e.tag, text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		e.delta = delta
		e.Name = strings.TrimSpace(match[2])
	}

	return nil
}

func (e TrackName) MarshalJSON() (encoded []byte, err error) {
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

func (e *TrackName) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Name  string    `json:"name"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagTrackName) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagTrackName
		e.Type = lib.TypeTrackName
		e.Name = t.Name
	}

	return nil
}
