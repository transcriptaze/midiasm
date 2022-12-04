package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type ProgramName struct {
	event
	Name string
}

func MakeProgramName(tick uint64, delta lib.Delta, name string, bytes ...byte) ProgramName {
	return ProgramName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   lib.TypeProgramName,
		},
		Name: name,
	}
}

func (e *ProgramName) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	name := string(data)

	*e = MakeProgramName(tick, delta, name, bytes...)

	return nil
}

func (p ProgramName) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(p.Status),
		byte(p.Type),
		byte(len(p.Name)),
	},
		[]byte(p.Name)...), nil
}

func (p *ProgramName) UnmarshalText(bytes []byte) error {
	p.tick = 0
	p.delta = 0
	p.bytes = []byte{}
	p.tag = lib.TagProgramName
	p.Status = 0xff
	p.Type = lib.TypeProgramName

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ProgramName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid ProgramName event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		p.delta = delta
		p.Name = string(match[2])
	}

	return nil
}

func (e ProgramName) MarshalJSON() (encoded []byte, err error) {
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

func (e *ProgramName) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Name  string    `json:"name"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagProgramName) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagProgramName
		e.Type = lib.TypeProgramName
		e.Name = t.Name
	}

	return nil
}
