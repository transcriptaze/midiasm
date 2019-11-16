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

	events []event.Event
}

func (chunk *MTrk) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	events := make([]event.Event, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	err := error(nil)
	event := Chunk(nil)

	for err == nil {
		event, err = parse(r)
		if err == nil && event != nil {
			events = append(events, event)
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

	for _, event := range chunk.events {
		event.Render(w)
	}

	fmt.Fprintln(w)
}

func parse(r *bufio.Reader) (event.Event, error) {
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

	if b == 0xff {
		return metaevent.Parse(delta, b, bytes, r)
	} else if b == 0xf0 || b == 0xf7 {
		panic(fmt.Sprintf("NOT IMPLEMENTED: SYSEX EVENT @%d:  %02x\n", delta, b))
	} else {
		return midievent.Parse(delta, b, bytes, r)
	}

	return nil, nil
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

		l <<= 8
		l += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return l, bytes, nil
}
