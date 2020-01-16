package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type MidiEvent struct {
	events.Event
	Channel types.Channel
}

type Note struct {
	Value byte
	Name  string
}

func (e MidiEvent) String() string {
	return fmt.Sprintf("%s %v", e.Event, e.Status)
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

func (e *MidiEvent) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%-2v", e.String(), e.Tag, e.Channel)
}

func Parse(e events.Event, r io.ByteReader, ctx *context.Context) (events.IEvent, error) {
	event := MidiEvent{
		Event:   e,
		Channel: types.Channel((e.Status) & 0x0F),
	}

	rr := reader{r, &event}

	switch e.Status & 0xF0 {
	case 0x80:
		event.Tag = "NoteOff"
		return NewNoteOff(ctx, &event, rr)

	case 0x90:
		event.Tag = "NoteOn"
		return NewNoteOn(ctx, &event, rr)

	case 0xA0:
		event.Tag = "PolyphonicPressure"
		return NewPolyphonicPressure(&event, rr)

	case 0xB0:
		event.Tag = "Controller"
		return NewController(&event, rr)

	case 0xC0:
		event.Tag = "ProgramChange"
		return NewProgramChange(&event, rr)

	case 0xD0:
		event.Tag = "ChannelPressure"
		return NewChannelPressure(&event, rr)

	case 0xE0:
		event.Tag = "PitchBend"
		return NewPitchBend(&event, rr)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", e.Status&0xF0)
}