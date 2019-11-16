package midievent

import (
	"bufio"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
)

func Parse(delta uint32, status byte, data []byte, r *bufio.Reader) (event.Event, error) {
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
			delta:    delta,
			status:   status,
			channel:  status & 0x0f,
			note:     note,
			velocity: velocity,
			bytes:    append(data, note, velocity),
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
			delta:    delta,
			status:   status,
			channel:  status & 0x0f,
			note:     note,
			velocity: velocity,
			bytes:    append(data, note, velocity),
		}, nil
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02x", status)
}
