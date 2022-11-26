package metaevent

import (
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type ProgramName struct {
	event
	Name string
}

func MakeProgramName(tick uint64, delta lib.Delta, name string) ProgramName {
	return ProgramName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x08, byte(len(name))}, []byte(name)...),
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   lib.TypeProgramName,
		},
		Name: name,
	}
}

func UnmarshalProgramName(tick uint64, delta lib.Delta, bytes []byte) (*ProgramName, error) {
	name := string(bytes)
	event := MakeProgramName(tick, delta, name)

	return &event, nil
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
