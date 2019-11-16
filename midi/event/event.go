package event

import (
	"io"
)

type Event interface {
	Render(io.Writer)
}
