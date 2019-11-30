package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"github.com/twystd/midiasm/midi/meta-events"
	"github.com/twystd/midiasm/midi/midi-events"
	"io"
)

type MTrk struct {
	tag    string
	length uint32
	data   []byte
	bytes  []byte

	Events []event.IEvent
}

func (chunk *MTrk) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	events := make([]event.IEvent, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	tick := uint32(0)
	err := error(nil)
	e := event.IEvent(nil)

	for err == nil {
		e, err = parse(r, tick)
		if err == nil && e != nil {
			tick += e.DeltaTime()
			events = append(events, e.(event.IEvent))
		}
	}

	if err != io.EOF {
		return err
	}

	chunk.tag = tag
	chunk.length = length
	chunk.data = data[8:]
	chunk.Events = events
	chunk.bytes = data

	return nil
}

func (chunk *MTrk) Render(w io.Writer) {
	context := event.Context{
		Scale: event.Sharps,
	}

	for _, b := range chunk.bytes[:8] {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "... ")
	fmt.Fprintf(w, "              %12s length:%-8d\n", chunk.tag, chunk.length)

	for _, e := range chunk.Events {
		e.Render(&context, w)
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

func parse(r *bufio.Reader, tick uint32) (event.IEvent, error) {
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

	e := event.Event{
		Tick:   uint64(tick + delta),
		Delta:  delta,
		Status: b,
		Bytes:  bytes,
	}

	if b == 0xff {
		return metaevent.Parse(e, r)
	} else if b == 0xf0 || b == 0xf7 {
		panic(fmt.Sprintf("NOT IMPLEMENTED: SYSEX EVENT @%d:  %02x\n", delta, b))
	} else {
		return midievent.Parse(e, r)
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
