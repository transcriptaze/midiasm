package event

import (
	"fmt"
	"io"
)

type IEvent interface {
	TickValue() uint64
	DeltaTime() uint32
	Render(w io.Writer)
}

type Event struct {
	Tick   uint64
	Delta  uint32
	Status byte
}

func (e *Event) TickValue() uint64 {
	return e.Tick
}

func (e *Event) DeltaTime() uint32 {
	return e.Delta
}

func (e Event) String() string {
	return fmt.Sprintf("tick:%-10d delta:%-10d", e.Tick, e.Delta)
}
