package assemble

import (
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

	bytes, err := assembler.Assemble(text)
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(bytes, smf) {
		t.Errorf("incorrectly assembled text file\n   expected:%+v\n   got:     %+v", smf, bytes)
	}

}
