package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/events/midi"
	"github.com/twystd/midiasm/midi/events/sysex"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type MTrk struct {
	Tag         string
	TrackNumber types.TrackNumber
	Length      uint32
	Bytes       types.Hex

	Events []*events.EventW
}

type reader struct {
	rdr   io.ByteReader
	event *events.Event
}

func (r reader) ReadByte() (byte, error) {
	b, err := r.rdr.ReadByte()
	if err == nil {
		r.event.Bytes = append(r.event.Bytes, b)
	}

	return b, err
}

func (chunk *MTrk) UnmarshalBinary(ctx *context.Context, data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	eventlist := make([]*events.EventW, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	tick := uint32(0)
	err := error(nil)
	var e *events.EventW = nil

	for err == nil {
		e, err = parse(r, tick, ctx)
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

func parse(r *bufio.Reader, tick uint32, ctx *context.Context) (*events.EventW, error) {
	bytes := make([]byte, 0)

	delta, m, err := vlq(r)
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, m...)

	bb, err := r.Peek(1)
	if err != nil {
		return nil, err
	}
	b := bb[0]

	// ... meta event
	if b == 0xff {
		ctx.RunningStatus = 0x00

		r.ReadByte()

		e := events.Event{
			Status: types.Status(b),
			Bytes:  append(bytes, b),
		}

		x, err := metaevent.Parse(&e, reader{r, &e}, ctx)
		return &events.EventW{
			Tick:  types.Tick(tick + delta),
			Delta: types.Delta(delta),
			Event: x,
		}, err
	}

	// ... SysEx event
	if b == 0xf0 || b == 0xf7 {
		ctx.RunningStatus = 0x00

		r.ReadByte()

		e := events.Event{
			Status: types.Status(b),
			Bytes:  append(bytes, b),
		}

		x, err := sysex.Parse(&e, reader{r, &e}, ctx)
		return &events.EventW{
			Tick:  types.Tick(tick + delta),
			Delta: types.Delta(delta),
			Event: x,
		}, err
	}

	// ... MIDI event
	if b < 0x80 && ctx.RunningStatus < 0x80 {
		return nil, fmt.Errorf("Unrecognised MIDI event: %02X", b&0xF0)
	}

	e := events.Event{
		Status: types.Status(b),
		Bytes:  bytes,
	}

	if b < 0x80 {
		e.Status = ctx.RunningStatus
	} else {
		r.ReadByte()
		e.Bytes = append(bytes, b)
	}

	ctx.RunningStatus = e.Status

	x, err := midievent.Parse(&e, reader{r, &e}, ctx)
	return &events.EventW{
		Tick:  types.Tick(tick + delta),
		Delta: types.Delta(delta),
		Event: x,
	}, err
}

func vlq(r *bufio.Reader) (uint32, []byte, error) {
	l := uint32(0)
	bytes := make([]byte, 0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		bytes = append(bytes, b)

		l <<= 7
		l += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return l, bytes, nil
}
