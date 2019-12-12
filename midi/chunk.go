package midi

import (
	"io"
)

type Chunk interface {
	Print(io.Writer)
}
