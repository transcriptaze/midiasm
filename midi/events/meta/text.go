package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type Text struct {
	event
	Text string
}

func MakeText(tick uint64, delta uint32, text string) Text {
	return Text{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x01, byte(len(text))}, []byte(text)...),
			tag:    lib.TagText,
			Status: 0xff,
			Type:   lib.TypeText,
		},
		Text: text,
	}
}

func UnmarshalText(tick uint64, delta uint32, bytes []byte) (*Text, error) {
	text := string(bytes)
	event := MakeText(tick, delta, text)

	return &event, nil
}

func (e Text) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.Text)),
	},
		[]byte(e.Text)...), nil
}

func (e *Text) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Text\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Text event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		e.tick = 0
		e.delta = uint32(delta)
		e.bytes = []byte{}
		e.tag = lib.TagText
		e.Status = 0xff
		e.Type = lib.TypeText
		e.Text = string(match[2])
	}

	return nil
}

func (e Text) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string `json:"tag"`
		Delta  uint32 `json:"delta"`
		Status byte   `json:"status"`
		Type   byte   `json:"type"`
		Text   string `json:"text"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Text:   e.Text,
	}

	return json.Marshal(t)
}

func (e *Text) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string `json:"tag"`
		Delta uint32 `json:"delta"`
		Text  string `json:"text"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagText) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagText
		e.Type = lib.TypeText
		e.Text = t.Text
	}

	return nil
}
