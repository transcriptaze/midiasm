package midievent

import (
	"bufio"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
)

type MidiEvent struct {
	event.Event
	channel byte
	bytes   []byte
}

func (e *MidiEvent) DeltaTime() uint32 {
	return e.Delta
}

func Parse(event event.Event, data []byte, r *bufio.Reader) (event.IEvent, error) {
	status := event.Status & 0xf0
	channel := event.Status & 0x0f

	if status == 0x80 {
		note, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		velocity, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		return &NoteOff{
			MidiEvent: MidiEvent{
				Event:   event,
				channel: channel,
				bytes:   append(data, note, velocity),
			},
			note:     note,
			velocity: velocity,
		}, nil
	}

	if status == 0x90 {
		note, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		velocity, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		return &NoteOn{
			MidiEvent: MidiEvent{
				Event:   event,
				channel: status,
				bytes:   append(data, note, velocity),
			},
			note:     note,
			velocity: velocity,
		}, nil
	}

	if status == 0xB0 {
		controller, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		value, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		return &Controller{
			MidiEvent: MidiEvent{
				Event:   event,
				channel: channel,
				bytes:   append(data, controller, value),
			},
			controller: controller,
			value:      value,
		}, nil
	}

	if status == 0xC0 {
		program, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		return &ProgramChange{
			MidiEvent: MidiEvent{
				Event:   event,
				channel: channel,
				bytes:   append(data, program),
			},
			program: program,
		}, nil
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", status)
}
