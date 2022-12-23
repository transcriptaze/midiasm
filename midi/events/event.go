package events

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TEvent interface {
	metaevent.TMetaEvent | midievent.TMidiEvent | sysex.TSysExEvent

	Tick() uint64
	Delta() uint32
	Bytes() []byte
}

type IEvent interface {
	Tick() uint64
	Delta() uint32
	Bytes() []byte
}

type Event struct {
	Event IEvent `json:"event"`
}

func NewEvent(e any, bytes ...byte) *Event {
	if v, ok := e.(IEvent); ok {
		return &Event{
			Event: v,
		}
	}

	panic(fmt.Sprintf("Invalid event (%v) - missing 'tick'", e))
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

func (e Event) Bytes() lib.Hex {
	if v, ok := e.Event.(IEvent); ok {
		return v.Bytes()
	}

	panic(fmt.Sprintf("Invalid event (%v) - missing 'bytes'", e))
}

func Is[E TEvent](e Event) bool {
	p := new(E)

	u := reflect.TypeOf(*p)
	v := reflect.TypeOf(e.Event)

	return u == v
}

func IsTrack0Event(e *Event) bool {
	switch e.Event.(type) {
	case
		metaevent.Tempo,
		metaevent.TimeSignature,
		metaevent.TrackName,
		metaevent.SMPTEOffset,
		metaevent.Copyright,
		metaevent.EndOfTrack:
		return true

	default:
		return false
	}
}

func IsTrack1Event(e *Event) bool {
	switch e.Event.(type) {
	case
		metaevent.Tempo,
		metaevent.SMPTEOffset:
		return false

	default:
		return true
	}
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
