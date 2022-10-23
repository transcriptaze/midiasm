package midievent

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type Controller struct {
	event
	Controller types.Controller
	Value      byte
}

func NewController(ctx *context.Context, tick uint64, delta uint32, r io.ByteReader, status types.Status) (*Controller, error) {
	if status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%v): expected 'Bx'", status)
	}

	channel := uint8(status & 0x0f)

	controller, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	value, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if ctx != nil && controller == 0x00 {
		ctx.ProgramBank[channel] = (ctx.ProgramBank[channel] & 0x003f) | ((uint16(value) & 0x003f) << 7)
	}

	if ctx != nil && controller == 0x20 {
		ctx.ProgramBank[channel] = (ctx.ProgramBank[channel] & (0x003f << 7)) | (uint16(value) & 0x003f)
	}

	return &Controller{
		event: event{
			tick:  tick,
			delta: types.Delta(delta),
			bytes: []byte{0x00, byte(status), controller, value},

			tag:     types.TagController,
			Status:  status,
			Channel: types.Channel(channel),
		},
		Controller: types.LookupController(controller),
		Value:      value,
	}, nil
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
	c.tag = types.TagController

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Controller.*\s+channel:([0-9]+)\s+([0-9]+)(?:.*)?value:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid Controller event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if controller, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if value, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid Controller channel (%v)", channel)
	} else {
		c.delta = types.Delta(delta)
		c.bytes = []byte{0x00, byte(0xb0 | uint8(channel&0x0f)), uint8(controller), uint8(value)}
		c.Status = types.Status(0xb0 | uint8(channel&0x0f))
		c.Channel = types.Channel(channel)
		c.Controller = types.LookupController(uint8(controller))
		c.Value = uint8(value)
	}

	return nil
}
