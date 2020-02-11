package midievent

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Note struct {
	Value byte
	Name  string
	Alias string
}

type reader struct {
	pushed io.ByteReader
	rdr    io.ByteReader
	event  *events.Event
}

func (r reader) ReadByte() (byte, error) {
	if r.pushed != nil {
		if b, err := r.pushed.ReadByte(); err == nil {
			return b, nil
		} else if err != io.EOF {
			return b, err
		}
	}

	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.Bytes = append(r.event.Bytes, b)
	}

	return b, err
}

func Parse(event events.Event, r io.ByteReader, ctx *context.Context) (interface{}, error) {
	rr := reader{nil, r, &event}

	//FIXME Ewwwwww :-(
	if event.Status < 0x80 && !ctx.HasRunningStatus() {
		return nil, fmt.Errorf("Unrecognised MIDI event: %02X", event.Status&0xF0)
	} else if event.Status < 0x80 {
		rr.pushed = bytes.NewReader([]byte{byte(event.Status)})
		event.Status = types.Status(ctx.GetRunningStatus())
	}

	ctx.PutRunningStatus(byte(event.Status))

	switch event.Status & 0xF0 {
	case 0x80:
		return NewNoteOff(ctx, &event, rr)

	case 0x90:
		return NewNoteOn(ctx, &event, rr)

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

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", event.Status&0xF0)
}
