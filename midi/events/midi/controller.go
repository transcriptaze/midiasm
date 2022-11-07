package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type Controller struct {
	event
	Controller lib.Controller
	Value      byte
}

func MakeController(tick uint64, delta uint32, channel lib.Channel, controller lib.Controller, value byte, bytes ...byte) Controller {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	return Controller{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagController,
			Status:  or(0xb0, channel),
			Channel: channel,
		},
		Controller: controller,
		Value:      value,
	}
}

func UnmarshalController(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (*Controller, error) {
	if status&0xf0 != 0xb0 {
		return nil, fmt.Errorf("Invalid Controller status (%v): expected 'Bx'", status)
	}

	var channel = lib.Channel(status & 0x0f)
	var controller byte
	var value byte

	if c, err := r.ReadByte(); err != nil {
		return nil, err
	} else {
		controller = c
	}

	if v, err := r.ReadByte(); err != nil {
		return nil, err
	} else {
		value = v
	}

	if ctx != nil && controller == 0x00 {
		c := uint8(channel)
		ctx.ProgramBank[c] = (ctx.ProgramBank[c] & 0x003f) | ((uint16(value) & 0x003f) << 7)
	}

	if ctx != nil && controller == 0x20 {
		c := uint8(channel)
		ctx.ProgramBank[c] = (ctx.ProgramBank[c] & (0x003f << 7)) | (uint16(value) & 0x003f)
	}

	event := MakeController(tick, delta, channel, lib.LookupController(controller), value, r.Bytes()...)

	return &event, nil
}

func (c Controller) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xb0 | c.Channel),
		c.Controller.ID,
		c.Value,
	}

	return
}

func (c *Controller) UnmarshalText(bytes []byte) error {
	c.tick = 0
	c.delta = 0
	c.bytes = []byte{}
	c.tag = lib.TagController
	c.Status = 0xb0

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Controller.*\s+channel:([0-9]+)\s+([0-9]+)(?:.*)?value:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid Controller event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if controller, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if value, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else {
		c.delta = delta
		c.Status = or(c.Status, channel)
		c.Channel = channel
		c.Controller = lib.LookupController(uint8(controller))
		c.Value = uint8(value)
	}

	return nil
}
