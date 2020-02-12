package events

import (
	"github.com/twystd/midiasm/midi/types"
)

type Event struct {
	Status types.Status
}

type EventW struct {
	Tick  types.Tick
	Delta types.Delta
	Event interface{}
	Bytes types.Hex
}
