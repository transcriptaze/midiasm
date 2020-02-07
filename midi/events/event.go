package events

import (
	"github.com/twystd/midiasm/midi/types"
)

type Event struct {
	Status types.Status
	Bytes  types.Hex
}

type EventW struct {
	Tick  types.Tick
	Delta types.Delta
	Event interface{}
}
