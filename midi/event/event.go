package event

import (
	"io"
)

type IEvent interface {
	DeltaTime() uint32
	Render(w io.Writer)
}
