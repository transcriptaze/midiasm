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
	"github.com/transcriptaze/midiasm/midi/io"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type MTrk struct {
	Tag         string
	TrackNumber lib.TrackNumber
	Length      uint32
	Bytes       lib.Hex `json:"-"`

	Events []*events.Event

	Context *context.Context
}

func NewMTrk() (*MTrk, error) {
	mtrk := MTrk{
		Tag:    "MTrk",
		Length: 0,
		Events: []*events.Event{},
		Bytes:  []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x00},
	}

	return &mtrk, nil
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
			tick += e.Delta()
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
			event.Mutate(-2, byte(v.Accidentals))

		case *midievent.NoteOn:
			v.Transpose(t.Context, steps)
			event.Mutate(-2, v.Note.Value)

		case *midievent.NoteOff:
			v.Transpose(t.Context, steps)
			event.Mutate(-2, v.Note.Value)
		}
	}
}

func parse(r *bufio.Reader, tick uint32, ctx *context.Context) (*events.Event, error) {
	// var buffer bytes.Buffer

	rr := IO.NewReader(r)

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

		e, err := metaevent.Parse(ctx, rr, uint64(tick)+uint64(delta), delta)

		return events.NewEvent(uint64(tick)+uint64(delta), delta, e, rr.Bytes()), err
	}

	// ... SysEx event
	if b == 0xf0 || b == 0xf7 {
		ctx.RunningStatus = 0x00

		rr.ReadByte()

		e, err := sysex.Parse(uint64(tick)+uint64(delta), delta, rr, lib.Status(b), ctx)

		return events.NewEvent(uint64(tick)+uint64(delta), delta, e, rr.Bytes()), err
	}

	// ... MIDI event
	if b < 0x80 && ctx.RunningStatus < 0x80 {
		return nil, fmt.Errorf("Unrecognised MIDI event: %02X", b&0xF0)
	}

	status := lib.Status(b)

	if b < 0x80 {
		status = ctx.RunningStatus
	} else {
		rr.ReadByte()
	}

	ctx.RunningStatus = status

	e, err := midievent.Parse(uint64(tick)+uint64(delta), delta, rr, status, ctx)

	return events.NewEvent(uint64(tick)+uint64(delta), delta, e, rr.Bytes()), err
}
