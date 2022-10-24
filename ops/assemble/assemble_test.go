package assemble

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"reflect"
	"testing"
)

//go:embed test-files/example.txt
var text []byte

//go:embed test-files/example.mid
var smf []byte

func TestTextAssemble(t *testing.T) {
	assembler := TextAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(text))
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smf) {
		t.Errorf("incorrectly assembled text file\nexpected:\n%+v\ngot:\n%+v", hex.Dump(smf), hex.Dump(encoded))
	}
}
