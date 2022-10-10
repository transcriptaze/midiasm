package assemble

import (
	"fmt"
	"io"
)

type JSONAssembler struct {
}

func NewJSONAssembler() JSONAssembler {
	return JSONAssembler{}
}

func (a JSONAssembler) Assemble(r io.Reader) ([]byte, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
