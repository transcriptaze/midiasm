package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type Text struct {
	event
	Text string
}

func MakeText(tick uint64, delta lib.Delta, text string, bytes ...byte) Text {
	return Text{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagText,
			Status: 0xff,
			Type:   lib.TypeText,
		},
		Text: text,
	}
}

func (e *Text) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	text := string(data)
	event := MakeText(tick, delta, text, bytes...)

	*e = event

	return nil
}

func (e Text) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.Text)),
	},
		[]byte(e.Text)...), nil
}

func (e *Text) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagText, remaining[0])
	} else if !equals(remaining[1], lib.TypeText) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagText, remaining[1])
	} else if text, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeText(0, delta, string(text), bytes...)
	}

	return nil
}

func (e *Text) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Text\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Text event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeText(0, delta, match[2], []byte{}...)
	}

	return nil
}

func (e Text) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
		Text   string    `json:"text"`
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
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Text  string    `json:"text"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagText) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeText(0, t.Delta, t.Text, []byte{}...)
	}

	return nil
}
