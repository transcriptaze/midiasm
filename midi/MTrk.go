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

	events []event.IEvent
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
	chunk.events = events
	chunk.bytes = data

	return nil
}

func (chunk *MTrk) Render(w io.Writer) {
	for _, b := range chunk.bytes[:8] {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "... ")
	fmt.Fprintf(w, "              %12s length:%-8d\n", chunk.tag, chunk.length)

	for _, e := range chunk.events {
		e.Render(w)
	}

	fmt.Fprintln(w)
}

func (chunk *MTrk) Notes(ppqn uint16,
	tx []struct {
		t     uint32
		tempo uint32
	}, w io.Writer) {

	var t float64 = 0.0
	var tick uint32 = 0
	var tempo float64 = float64(tx[0].tempo) / 1000000.0

	for _, e := range chunk.events {
		dt := e.DeltaTime()
		t += tempo * float64(dt) / float64(ppqn)
		tick += dt
		beat := float64(tick) / float64(ppqn)

		switch e.(type) {
		case *midievent.NoteOn:
			fmt.Fprintf(w, "NOTE ON  %-6d %.5f  %.5f\n", tick, beat, t)
		case *midievent.NoteOff:
			fmt.Fprintf(w, "NOTE OFF %-6d %.5f  %.5f\n", tick, beat, t)
		}
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
		Tick:   tick + delta,
		Delta:  delta,
		Status: b,
	}

	if b == 0xff {
		return metaevent.Parse(e, bytes, r)
	} else if b == 0xf0 || b == 0xf7 {
		panic(fmt.Sprintf("NOT IMPLEMENTED: SYSEX EVENT @%d:  %02x\n", delta, b))
	} else {
		return midievent.Parse(e, bytes, r)
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
