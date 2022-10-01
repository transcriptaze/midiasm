package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/types"
)

type MTrk struct {
	Tag         string
	TrackNumber types.TrackNumber
	Length      uint32
	Bytes       types.Hex `json:"-"`

	Events []*events.Event

	Context *context.Context
}

func (chunk *MTrk) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	eventlist := make([]*events.Event, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	tick := uint32(0)
	err := error(nil)
	var e *events.Event = nil

	for err == nil {
		e, err = parse(r, tick, chunk.Context)
		if err == nil && e != nil {
			tick += uint32(e.Delta)
			eventlist = append(eventlist, e)
		}
	}

	if err != io.EOF {
		return err
	}

	chunk.Tag = tag
	chunk.Length = length
	chunk.Events = eventlist
	chunk.Bytes = data

	return nil
}

func (t *MTrk) Transpose(steps int) {
	for _, event := range t.Events {
		switch v := event.Event.(type) {
		case *metaevent.KeySignature:
			v.Transpose(t.Context, steps)
			event.Bytes[len(event.Bytes)-2] = byte(v.Accidentals)

		case *midievent.NoteOn:
			v.Transpose(t.Context, steps)
			event.Bytes[len(event.Bytes)-2] = v.Note.Value

		case *midievent.NoteOff:
			v.Transpose(t.Context, steps)
			event.Bytes[len(event.Bytes)-2] = v.Note.Value
		}
	}
}

func parse(r *bufio.Reader, tick uint32, ctx *context.Context) (*events.Event, error) {
	var buffer bytes.Buffer

	rr := reader{r, &buffer}

	delta, err := events.VLQ(rr)
	if err != nil {
		return nil, err
	}

	bb, err := rr.Peek(1)
	if err != nil {
		return nil, err
	}
	b := bb[0]

	// ... meta event
	if b == 0xff {
		ctx.RunningStatus = 0x00

		rr.ReadByte()

		e, err := metaevent.Parse(ctx, rr, types.Status(b))

		return &events.Event{
			Tick:  types.Tick(tick + delta),
			Delta: types.Delta(delta),
			Bytes: buffer.Bytes(),
			Event: e,
		}, err
	}

	// ... SysEx event
	if b == 0xf0 || b == 0xf7 {
		ctx.RunningStatus = 0x00

		rr.ReadByte()

		e, err := sysex.Parse(rr, types.Status(b), ctx)

		return &events.Event{
			Tick:  types.Tick(tick + delta),
			Delta: types.Delta(delta),
			Bytes: buffer.Bytes(),
			Event: e,
		}, err
	}

	// ... MIDI event
	if b < 0x80 && ctx.RunningStatus < 0x80 {
		return nil, fmt.Errorf("Unrecognised MIDI event: %02X", b&0xF0)
	}

	status := types.Status(b)

	if b < 0x80 {
		status = ctx.RunningStatus
	} else {
		rr.ReadByte()
	}

	ctx.RunningStatus = status

	e, err := midievent.Parse(rr, status, ctx)

	return &events.Event{
		Tick:  types.Tick(tick + delta),
		Delta: types.Delta(delta),
		Bytes: buffer.Bytes(),
		Event: e,
	}, err
}
