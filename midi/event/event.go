package event

import (
	"bytes"
	"fmt"
	"io"
)

type IEvent interface {
	TickValue() uint64
	DeltaTime() uint32
	Render(*Context, io.Writer)
}

type Event struct {
	Tick   uint64
	Delta  uint32
	Status byte
	Bytes  []byte
}

func (e *Event) TickValue() uint64 {
	return e.Tick
}

func (e *Event) DeltaTime() uint32 {
	return e.Delta
}

func (e Event) String() string {
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "   ")

	for i := 5; i > len(e.Bytes); i-- {
		fmt.Fprintf(buffer, "   ")
	}

	for _, b := range e.Bytes {
		fmt.Fprintf(buffer, "%02X ", b)
	}

	if buffer.Len() > 42 {
		return fmt.Sprintf("%s\u2026  tick:%-10d delta:%-10d", buffer.String()[:41], e.Tick, e.Delta)
	} else {
		return fmt.Sprintf("%-42s  tick:%-10d delta:%-10d", buffer.String(), e.Tick, e.Delta)
	}
}
