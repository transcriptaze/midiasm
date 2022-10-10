package assemble

import (
	"bytes"
	_ "embed"
	"reflect"
	"testing"
)

//go:embed example.txt
var text []byte

//go:embed example.mid
var smf []byte

func TestTextAssemble(t *testing.T) {
	assembler := TextAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(text))
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smf) {
		t.Errorf("incorrectly assembled text file\n   expected:%+v\n   got:     %+v", smf, encoded)
	}

}
