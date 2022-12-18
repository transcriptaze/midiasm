package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type Controller struct {
	event
	Controller lib.Controller
	Value      byte
}

func MakeController(tick uint64, delta uint32, channel lib.Channel, controller lib.Controller, value byte, bytes ...byte) Controller {
	if channel > 15 {
		panic(fmt.Errorf("invalid channel (%v)", channel))
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

func (e *Controller) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0xb0 {
		return fmt.Errorf("Invalid %v status (%v): expected 'Bx'", lib.TagController, status)
	}

	if len(data) < 2 {
		return fmt.Errorf("Invalid %v data (%v): expected note and velocity", lib.TagController, data)
	}

	var channel = lib.Channel(status & 0x0f)
	var controller = data[0]
	var value = data[1]

	if channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	}

	event := MakeController(tick, delta, channel, lib.LookupController(controller), value, bytes...)

	if ctx != nil && event.Controller.ID == 0x00 {
		c := uint8(event.Channel)
		v := uint16(event.Value)
		ctx.ProgramBank[c] = (ctx.ProgramBank[c] & 0x003f) | ((v & 0x003f) << 7)
	}

	if ctx != nil && event.Controller.ID == 0x20 {
		c := uint8(event.Channel)
		v := uint16(event.Value)
		ctx.ProgramBank[c] = (ctx.ProgramBank[c] & (0x003f << 7)) | (v & 0x003f)
	}

	*e = event

	return nil
}

func (e Controller) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xb0 | e.Channel),
		e.Controller.ID,
		e.Value,
	}

	return
}

func (e *Controller) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 3 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeController.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagController, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 2 {
		return fmt.Errorf("Invalid Controller data")
	} else {
		controller := lib.LookupController(data[0])
		value := data[1]
		*e = MakeController(0, delta, lib.Channel(channel), controller, value, bytes...)
	}

	return nil
}

func (e *Controller) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Controller.*\s+channel:([0-9]+)\s+([0-9]+)(?:.*)?value:([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 5 {
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
		*e = MakeController(0, uint32(delta), lib.Channel(channel), lib.LookupController(uint8(controller)), uint8(value), []byte{}...)
	}

	return nil
}

func (e Controller) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag        string         `json:"tag"`
		Delta      lib.Delta      `json:"delta"`
		Status     byte           `json:"status"`
		Channel    lib.Channel    `json:"channel"`
		Controller lib.Controller `json:"controller"`
		Value      byte           `json:"value"`
	}{
		Tag:        fmt.Sprintf("%v", e.tag),
		Delta:      e.delta,
		Status:     byte(e.Status),
		Channel:    e.Channel,
		Controller: e.Controller,
		Value:      e.Value,
	}

	return json.Marshal(t)
}

func (e *Controller) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag        string         `json:"tag"`
		Delta      lib.Delta      `json:"delta"`
		Channel    lib.Channel    `json:"channel"`
		Controller lib.Controller `json:"controller"`
		Value      byte           `json:"value"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagController) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeController(0, uint32(t.Delta), t.Channel, t.Controller, t.Value, []byte{}...)
	}

	return nil
}
