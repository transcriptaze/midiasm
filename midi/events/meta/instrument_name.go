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

func (n InstrumentName) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(n.Status),
		byte(n.Type),
		byte(len(n.Name)),
	},
		[]byte(n.Name)...), nil
}

func (e *InstrumentName) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)InstrumentName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid InstrumentName event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagInstrumentName
		e.Status = 0xff
		e.Type = lib.TypeInstrumentName
		e.Name = match[2]
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
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagInstrumentName
		e.Type = lib.TypeInstrumentName
		e.Name = t.Name
	}

	return nil
}
