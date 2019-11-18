package event

import (
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
