package midievent

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
	"strings"
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
		fmt.Fprintf(buffer, "%02X ", b)
	}

	//fmt.Fprintf(buffer, "                                     ")
	//fmt.Fprintf(buffer, "%s %02X", e.Event, e.Status)

	//	return buffer.String()

	fmt.Fprintf(buffer, "%s", strings.Repeat(" ", 60-buffer.Len()))

	return fmt.Sprintf("%s %s %02X", buffer.String()[:60], e.Event, e.Status)
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
