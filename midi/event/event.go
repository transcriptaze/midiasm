package event

import (
	"io"
)

type Event interface {
	DeltaTime() uint32
	Render(w io.Writer)
}
