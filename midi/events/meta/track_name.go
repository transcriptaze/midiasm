package metaevent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type TrackName struct {
	event
	Name string
}

func MakeTrackName(tick uint64, delta uint32, name string) TrackName {
	n := lib.VLF(name)
	v, _ := n.MarshalBinary()

	return TrackName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x03}, v...),
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
		},
		Name: name,
	}
}

func UnmarshalTrackName(tick uint64, delta uint32, bytes []byte) (*TrackName, error) {
	name := string(bytes)
	event := MakeTrackName(tick, delta, name)

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
	} else {
		if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
			return err
		} else {
			e.delta = uint32(delta)
			e.Name = strings.TrimSpace(match[2])
		}
	}

	return nil
}

func (e TrackName) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string `json:"tag"`
		Delta  uint32 `json:"delta"`
		Status byte   `json:"status"`
		Type   byte   `json:"type"`
		Name   string `json:"name"`
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
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = lib.TagTrackName
	e.Type = lib.TypeTrackName

	t := struct {
		Tag   string `json:"tag"`
		Delta uint32 `json:"delta"`
		Name  string `json:"name"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != "TrackName" {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.delta = t.Delta
		e.Name = t.Name
	}

	return nil
}
