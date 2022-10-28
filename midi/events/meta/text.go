package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
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
			tag:    types.TagText,
			Status: 0xff,
			Type:   types.TypeText,
		},
		Text: text,
	}
}

func UnmarshalText(tick uint64, delta uint32, bytes []byte) (*Text, error) {
	text := string(bytes)
	event := MakeText(tick, delta, text)

	return &event, nil
}

func (t Text) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(t.Status),
		byte(t.Type),
		byte(len(t.Text)),
	},
		[]byte(t.Text)...), nil
}

func (t *Text) UnmarshalText(bytes []byte) error {
	t.tick = 0
	t.delta = 0
	t.bytes = []byte{}
	t.tag = types.TagText
	t.Status = 0xff
	t.Type = types.TypeText

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Text\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Text event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		t.delta = uint32(delta)
		t.Text = string(match[2])
	}

	return nil
}
