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
