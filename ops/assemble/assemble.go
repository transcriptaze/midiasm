package assemble

type Assembler interface {
	Assemble([]byte) ([]byte, error)
}
