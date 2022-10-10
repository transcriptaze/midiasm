package assemble

import (
	"io"
)

type Assembler interface {
	Assemble(io.Reader) ([]byte, error)
}
