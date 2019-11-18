package midievent

import (
	"bufio"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
)

type MidiEvent struct {
	delta   uint32
	status  byte
	channel byte
	bytes   []byte
}

func (e *MidiEvent) DeltaTime() uint32 {
	return e.delta
}

func Parse(delta uint32, status byte, data []byte, r *bufio.Reader) (event.IEvent, error) {
	if status&0xf0 == 0x80 {
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
				delta:   delta,
				status:  status,
				channel: status & 0x0f,
				bytes:   append(data, note, velocity),
			},
			note:     note,
			velocity: velocity,
		}, nil
	}

	if status&0xf0 == 0x90 {
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
				delta:   delta,
				status:  status,
				channel: status & 0x0f,
				bytes:   append(data, note, velocity),
			},
			note:     note,
			velocity: velocity,
		}, nil
	}

	if status&0xf0 == 0xB0 {
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
				delta:   delta,
				status:  status,
				channel: status & 0x0f,
				bytes:   append(data, controller, value),
			},
			controller: controller,
			value:      value,
		}, nil
	}

	if status&0xf0 == 0xC0 {
		program, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		return &ProgramChange{
			MidiEvent: MidiEvent{
				delta:   delta,
				status:  status,
				channel: status & 0x0f,
				bytes:   append(data, program),
			},
			program: program,
		}, nil
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", status)
}
