package assemble

import ()

type Assembler interface {
	Assemble([]byte) ([]byte, error)
}
