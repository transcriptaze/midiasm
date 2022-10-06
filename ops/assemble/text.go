package assemble

import (
	"fmt"
)

type TextAssembler struct {
}

func NewTextAssembler() TextAssembler {
	return TextAssembler{}
}

func (a TextAssembler) Assemble([]byte) ([]byte, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
