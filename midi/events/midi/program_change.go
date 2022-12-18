package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type ProgramChange struct {
	event
	Bank    uint16
	Program byte
}

func MakeProgramChange(tick uint64, delta uint32, channel lib.Channel, bank uint16, program uint8, bytes ...byte) ProgramChange {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	return ProgramChange{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagProgramChange,
			Status:  or(0xc0, channel),
			Channel: channel,
		},
		Bank:    bank,
		Program: program,
	}
}

func (e *ProgramChange) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0xc0 {
		return fmt.Errorf("Invalid %v status (%v): expected 'Cx'", lib.TagProgramChange, status)
	}

	if len(data) < 1 {
		return fmt.Errorf("Invalid %v data (%v): expected note and velocity", lib.TagProgramChange, data)
	}

	var channel = lib.Channel(status & 0x0f)
	var bank uint16
	var program = data[0]

	if ctx != nil {
		bank = ctx.ProgramBank[uint8(channel)]
	}

	*e = MakeProgramChange(tick, delta, channel, bank, program, bytes...)

	return nil
}

func (e ProgramChange) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xc0 | e.Channel),
		byte(e.Program),
	}

	return
}

func (e *ProgramChange) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeProgramChange.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagProgramChange, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 1 {
		return fmt.Errorf("Invalid ProgramChange data")
	} else {
		bank := uint16(0)
		program := data[0]

		*e = MakeProgramChange(0, delta, lib.Channel(channel), bank, program, bytes...)
	}

	return nil
}

func (e *ProgramChange) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ProgramChange\s+channel:([0-9]+)\s+bank:([0-9]+),\s*program:([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid ProgramChange event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if bank, err := strconv.ParseUint(match[3], 10, 16); err != nil {
		return err
	} else if program, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else {
		*e = MakeProgramChange(0, uint32(delta), lib.Channel(channel), uint16(bank), uint8(program), []byte{}...)
	}

	return nil
}

func (e ProgramChange) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag     string      `json:"tag"`
		Delta   lib.Delta   `json:"delta"`
		Status  byte        `json:"status"`
		Channel lib.Channel `json:"channel"`
		Bank    uint16      `json:"bank"`
		Program uint8       `json:"program"`
	}{
		Tag:     fmt.Sprintf("%v", e.tag),
		Delta:   e.delta,
		Status:  byte(e.Status),
		Channel: e.Channel,
		Bank:    e.Bank,
		Program: e.Program,
	}

	return json.Marshal(t)
}

func (e *ProgramChange) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag     string      `json:"tag"`
		Delta   lib.Delta   `json:"delta"`
		Channel lib.Channel `json:"channel"`
		Bank    uint16      `json:"bank"`
		Program uint8       `json:"program"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagProgramChange) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeProgramChange(0, uint32(t.Delta), t.Channel, t.Bank, t.Program, []byte{}...)
	}

	return nil
}
