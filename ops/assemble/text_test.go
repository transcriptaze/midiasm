package assemble

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"reflect"
	"testing"
)

//go:embed test-files/reference.txt
var reference []byte

//go:embed test-files/reference.mid
var smf []byte

func TestTextReference(t *testing.T) {
	assembler := TextAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(reference))
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smf) {
		t.Errorf("incorrectly assembled text file\nexpected:\n%+v\ngot:\n%+v", hex.Dump(smf), hex.Dump(encoded))
	}
}
