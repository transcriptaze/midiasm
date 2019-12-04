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
	"io"
)

type MTrk struct {
	tag    string
	length uint32
	data   []byte
	bytes  []byte

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

	for err == nil {
		e, err = parse(r, tick)
		if err == nil && e != nil {
			tick += e.DeltaTime()
			eventlist = append(eventlist, e.(events.IEvent))
		}
	}

	if err != io.EOF {
		return err
	}

	chunk.tag = tag
	chunk.length = length
	chunk.data = data[8:]
	chunk.Events = eventlist
	chunk.bytes = data

	return nil
}

func (chunk *MTrk) Render(w io.Writer) {
	context := context.Context{
		Scale: context.Sharps,
	}

	buffer := new(bytes.Buffer)
	for i, b := range chunk.bytes[:8] {
		if i == 0 {
			fmt.Fprintf(buffer, "%02X", b)
		} else {
			fmt.Fprintf(buffer, " %02X", b)
		}
	}

	if len(chunk.bytes) > 8 {
		fmt.Fprintf(buffer, "\u2026")
	} else {
		fmt.Fprintf(buffer, " ")
	}

	fmt.Fprintf(w, "%s            %12s length:%-8d\n", buffer.String(), chunk.tag, chunk.length)

	for _, e := range chunk.Events {
		e.Render(&context, w)
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

func parse(r *bufio.Reader, tick uint32) (events.IEvent, error) {
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
