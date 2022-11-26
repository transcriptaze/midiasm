package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/lib"
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

func UnmarshalController(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (*Controller, error) {
	if status&0xf0 != 0xb0 {
		return nil, fmt.Errorf("Invalid %v status (%v): expected 'Bx'", lib.TagController, status)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid %v data (%v): expected note and velocity", lib.TagController, data)
	}

	var channel = lib.Channel(status & 0x0f)
	var controller = data[0]
	var value = data[1]

	event := MakeController(tick, delta, channel, lib.LookupController(controller), value, bytes...)

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
