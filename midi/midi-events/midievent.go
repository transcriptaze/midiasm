package midievent

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
)

type MidiEvent struct {
	event.Event
	Channel byte
	bytes   []byte
}

func (e MidiEvent) String() string {
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "   ")

	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(buffer, "   ")
	}

	for _, b := range e.bytes {
		fmt.Fprintf(buffer, "%02x ", b)
	}

	fmt.Fprintf(buffer, "                                     ")
	fmt.Fprintf(buffer, "%s %02X", e.Event, e.Status)

	return buffer.String()
}

func Parse(event event.Event, data []byte, r *bufio.Reader) (event.IEvent, error) {
	midiEvent := MidiEvent{
		Event:   event,
		Channel: event.Status & 0x0F,
		bytes:   data,
	}

	switch event.Status & 0xF0 {
	case 0x80:
		return NewNoteOff(midiEvent, r)

	case 0x90:
		return NewNoteOn(midiEvent, r)

	case 0xA0:
		return NewPolyphonicPressure(midiEvent, r)

	case 0xB0:
		return NewController(midiEvent, r)

	case 0xC0:
		return NewProgramChange(midiEvent, r)

	case 0xD0:
		return NewChannelPressure(midiEvent, r)

	case 0xE0:
		return NewPitchBend(midiEvent, r)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", event.Status&0xF0)
}
