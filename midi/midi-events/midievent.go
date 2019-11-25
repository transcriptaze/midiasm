package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type MidiEvent struct {
	event.Event
	Channel byte
}

func (e MidiEvent) String() string {
	return fmt.Sprintf("%s %02X", e.Event, e.Status)
}

type reader struct {
	rdr   io.ByteReader
	event *MidiEvent
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.Bytes = append(r.event.Bytes, b)
	}

	return b, err
}

func Parse(e event.Event, r io.ByteReader) (event.IEvent, error) {
	event := MidiEvent{
		Event:   e,
		Channel: e.Status & 0x0F,
	}

	rr := reader{r, &event}

	switch e.Status & 0xF0 {
	case 0x80:
		return NewNoteOff(&event, rr)

	case 0x90:
		return NewNoteOn(&event, rr)

	case 0xA0:
		return NewPolyphonicPressure(&event, rr)

	case 0xB0:
		return NewController(&event, rr)

	case 0xC0:
		return NewProgramChange(&event, rr)

	case 0xD0:
		return NewChannelPressure(&event, rr)

	case 0xE0:
		return NewPitchBend(&event, rr)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", e.Status&0xF0)
}
