package event

import (
	"fmt"
	"io"
)

type IEvent interface {
	DeltaTime() uint32
	Render(w io.Writer)
}

type Event struct {
	Tick   uint32
	Delta  uint32
	Status byte
}

func (e Event) String() string {
	return fmt.Sprintf("tick:%-10d delta:%-10d", e.Tick, e.Delta)
}
