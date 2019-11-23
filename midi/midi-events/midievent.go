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

func Parse(event event.Event, data []byte, r io.ByteReader) (event.IEvent, error) {
	midiEvent := MidiEvent{
		Event:   event,
		Channel: event.Status & 0x0F,
		bytes:   data,
	}

	rr := reader{r, &midiEvent}

	switch event.Status & 0xF0 {
	case 0x80:
		return NewNoteOff(&midiEvent, rr)

	case 0x90:
		return NewNoteOn(&midiEvent, rr)

	case 0xA0:
		return NewPolyphonicPressure(&midiEvent, rr)

	case 0xB0:
		return NewController(&midiEvent, rr)

	case 0xC0:
		return NewProgramChange(&midiEvent, rr)

	case 0xD0:
		return NewChannelPressure(&midiEvent, rr)

	case 0xE0:
		return NewPitchBend(&midiEvent, rr)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", event.Status&0xF0)
}
