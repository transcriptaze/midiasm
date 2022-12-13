package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type Copyright struct {
	event
	Copyright string
}

func MakeCopyright(tick uint64, delta lib.Delta, copyright string, bytes ...byte) Copyright {
	return Copyright{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
		},
		Copyright: copyright,
	}
}

func (e *Copyright) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	copyright := string(data)
	*e = MakeCopyright(tick, delta, copyright, bytes...)

	return nil
}

func (e Copyright) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.Copyright)),
	},
		[]byte(e.Copyright)...), nil
}

func (e *Copyright) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagCopyright, remaining[0])
	} else if !equals(remaining[1], lib.TypeCopyright) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagCopyright, remaining[1])
	} else if copyright, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeCopyright(0, delta, string(copyright), bytes...)
	}

	return nil
}

func (e *Copyright) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Copyright\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Copyright event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeCopyright(0, delta, match[2], []byte{}...)
	}

	return nil
}

func (e Copyright) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag       string    `json:"tag"`
		Delta     lib.Delta `json:"delta"`
		Status    byte      `json:"status"`
		Type      byte      `json:"type"`
		Copyright string    `json:"copyright"`
	}{
		Tag:       fmt.Sprintf("%v", e.tag),
		Delta:     e.delta,
		Status:    byte(e.Status),
		Type:      byte(e.Type),
		Copyright: e.Copyright,
	}

	return json.Marshal(t)
}

func (e *Copyright) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag       string    `json:"tag"`
		Delta     lib.Delta `json:"delta"`
		Copyright string    `json:"copyright"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagCopyright) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeCopyright(0, t.Delta, t.Copyright, []byte{}...)
	}

	return nil
}
