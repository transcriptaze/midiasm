package events

import (
	"github.com/twystd/midiasm/midi/types"
)

type Event struct {
	Tick  types.Tick
	Delta types.Delta
	Event interface{}
	Bytes types.Hex
}

type EventReader interface {
	ReadByte() (byte, error)
	Peek(n int) ([]byte, error)
	ReadVLF() ([]byte, error)
	VLQ() (uint32, error)
}
