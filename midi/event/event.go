package event

import (
	"bytes"
	"fmt"
	"io"
	"strings"
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

	fmt.Fprintf(buffer, "%s", strings.Repeat(" ", 60-buffer.Len()))

	return fmt.Sprintf("%s tick:%-10d delta:%-10d", buffer.String()[:60], e.Tick, e.Delta)
}
