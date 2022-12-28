package metaevent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

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

func (e *TrackName) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	name := string(data)
	event := MakeTrackName(tick, delta, name, bytes...)

	*e = event

	return nil
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

func (e *TrackName) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagTrackName, remaining[0])
	} else if !equals(remaining[1], lib.TypeTrackName) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagTrackName, remaining[1])
	} else if name, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeTrackName(0, delta, string(name), bytes...)
	}

	return nil
}

func (e *TrackName) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)TrackName\s+(.*)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid %v event (%v)", e.tag, text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeTrackName(0, delta, match[2], []byte{}...)
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
		*e = MakeTrackName(0, t.Delta, t.Name, []byte{}...)
	}

	return nil
}
