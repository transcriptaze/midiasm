package assemble

import (
	"fmt"
)

type JSONAssembler struct {
}

func NewJSONAssembler() JSONAssembler {
	return JSONAssembler{}
}

func (a JSONAssembler) Assemble([]byte) ([]byte, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
