package midievent

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
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

type reader struct {
	rdr   io.ByteReader
	event *MidiEvent
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.bytes = append(r.event.bytes, b)
	}

	return b, err
}

func Parse(e event.Event, data []byte, r io.ByteReader) (event.IEvent, error) {
	event := MidiEvent{
		Event:   e,
		Channel: e.Status & 0x0F,
		bytes:   data,
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
