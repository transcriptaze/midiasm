package events

import (
	"fmt"
	"io"
	"strings"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/types"
)

type TEvent interface {
	metaevent.TMetaEvent | midievent.TMidiEvent | sysex.TSysExEvent
}

type IEvent interface {
	Tick() uint64
	Delta() uint32
	Bytes() []byte
}

// type Event[E TEvent] struct {
// 	Event E
// }

// func NewEvent[E TEvent](e E) *Event[E] {
// 	return &Event[E]{
// 		Event: e,
// 	}
// }

type Event struct {
	Event any
}

func NewEvent(e any) *Event {
	return &Event{
		Event: e,
	}
}

func (e Event) Tick() uint64 {
	if v, ok := e.Event.(IEvent); ok {
		return v.Tick()
	}

	panic(fmt.Sprintf("Invalid event (%v) - missing 'tick'", e))
}

func (e Event) Delta() uint32 {
	if v, ok := e.Event.(IEvent); ok {
		return v.Delta()
	}

	panic(fmt.Sprintf("Invalid event (%v) - missing 'delta'", e))
}

func (e Event) Bytes() types.Hex {
	if v, ok := e.Event.(IEvent); ok {
		return v.Bytes()
	}

	panic(fmt.Sprintf("Invalid event (%v) - missing 'bytes'", e))
}

func IsEndOfTrack(e *Event) bool {
	_, ok := e.Event.(*metaevent.EndOfTrack)

	return ok
}

func Clean(e any) string {
	t := ""

	if evt, ok := e.(*Event); ok {
		t = fmt.Sprintf("%T", evt.Event)
	} else {
		t = fmt.Sprintf("%T", e)
	}

	t = strings.TrimPrefix(t, "*")
	t = strings.TrimPrefix(t, "metaevent.")
	t = strings.TrimPrefix(t, "midievent.")
	t = strings.TrimPrefix(t, "sysex.")

	return t
}

func VLF(r io.ByteReader) ([]byte, error) {
	N, err := VLQ(r)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, N)

	for i := 0; i < int(N); i++ {
		if b, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			bytes[i] = b
		}
	}

	return bytes, nil
}

func VLQ(r io.ByteReader) (uint32, error) {
	vlq := uint32(0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		vlq <<= 7
		vlq += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return vlq, nil
}
