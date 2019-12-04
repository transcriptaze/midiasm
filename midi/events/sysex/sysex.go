package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"io"
)

type SysExEvent struct {
	events.Event
}

func (e SysExEvent) String() string {
	return fmt.Sprintf("%s %02X", e.Event, e.Status)
}

func Parse(e events.Event, r io.ByteReader) (events.IEvent, error) {
	if e.Status != 0xF0 && e.Status != 0xF7 {
		return nil, fmt.Errorf("Invalid SysEx tag (%02x): expected 'F0' or 'F7'", e.Status)
	}

	event := SysExEvent{
		Event: e,
	}

	switch event.Status {
	case 0xf0:
		return NewSysExSingleMessage(&event, r)

	case 0xf7:
		return NewSysExContinuationMessage(&event, r)
	}

	return nil, fmt.Errorf("Unrecognised SYSEX event: %02X", event.Status)
}

func read(r io.ByteReader) ([]byte, error) {
	N, err := vlq(r)
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

func vlq(r io.ByteReader) (uint32, error) {
	l := uint32(0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		l <<= 7
		l += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return l, nil
}
