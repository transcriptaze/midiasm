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
	"text/template"
)

type MTrk struct {
	Tag         string
	TrackNumber types.TrackNumber
	Length      uint32
	Bytes       types.Hex

	Events []events.IEvent
}

func (chunk *MTrk) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	eventlist := make([]events.IEvent, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	tick := uint32(0)
	err := error(nil)
	e := events.IEvent(nil)
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: false,
	}

	for err == nil {
		e, err = parse(r, tick, &ctx)
		if err == nil && e != nil {
			tick += e.DeltaTime()
			eventlist = append(eventlist, e.(events.IEvent))
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

func (chunk *MTrk) Print(w io.Writer) error {
	format := "{{slice .Bytes 0 8}}…                    {{.Tag}} {{.TrackNumber}} length:{{.Length}}"
	tmpl, err := template.New("MTrk").Parse(format)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, chunk)
	if err != nil {
		return err
	}

	for _, e := range chunk.Events {
		fmt.Fprintln(w)
		e.Render(w)
	}

	return nil
}

func parse(r *bufio.Reader, tick uint32, ctx *context.Context) (events.IEvent, error) {
	bytes := make([]byte, 0)

	delta, m, err := vlq(r)
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, m...)

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, b)

	e := events.Event{
		Tick:   types.Tick(tick + delta),
		Delta:  types.Delta(delta),
		Status: types.Status(b),
		Bytes:  bytes,
	}

	if b == 0xff {
		return metaevent.Parse(e, r, ctx)
	} else if b == 0xf0 || b == 0xf7 {
		return sysex.Parse(e, r, ctx)
	} else {
		return midievent.Parse(e, r, ctx)
	}
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
