package midi

import (
	"io"
)

type Chunk interface {
	Render(io.Writer)
}
