package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type InstrumentName struct {
	event
	Name string
}

func MakeInstrumentName(tick uint64, delta lib.Delta, name string, bytes ...byte) InstrumentName {
	return InstrumentName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
		},
		Name: name,
	}
}

func (e *InstrumentName) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	name := string(data)

	*e = MakeInstrumentName(tick, delta, name, bytes...)

	return nil
}

func (e InstrumentName) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.Name)),
	},
		[]byte(e.Name)...), nil
}

func (e *InstrumentName) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagInstrumentName, remaining[0])
	} else if !equals(remaining[1], lib.TypeInstrumentName) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagInstrumentName, remaining[1])
	} else if name, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeInstrumentName(0, delta, string(name), bytes...)
	}

	return nil
}

func (e *InstrumentName) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)InstrumentName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid InstrumentName event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeInstrumentName(0, delta, match[2], []byte{}...)
	}

	return nil
}

func (e InstrumentName) MarshalJSON() (encoded []byte, err error) {
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

func (e *InstrumentName) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Name  string    `json:"name"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagInstrumentName) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeInstrumentName(0, t.Delta, t.Name, []byte{}...)
	}

	return nil
}
