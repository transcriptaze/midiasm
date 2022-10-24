package assemble

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"reflect"
	"testing"
)

//go:embed test-files/example.txt
var example []byte

//go:embed test-files/example.mid
var smf1 []byte

//go:embed test-files/reference.txt
var reference []byte

//go:embed test-files/reference.mid
var smf2 []byte

func TestTextExample(t *testing.T) {
	assembler := TextAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(example))
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smf1) {
		t.Errorf("incorrectly assembled text file\nexpected:\n%+v\ngot:\n%+v", hex.Dump(smf1), hex.Dump(encoded))
	}
}

func TestTextReference(t *testing.T) {
	assembler := TextAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(reference))
	if err != nil {
		t.Fatalf("error assembler text file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smf2) {
		t.Errorf("incorrectly assembled text file\nexpected:\n%+v\ngot:\n%+v", hex.Dump(smf2), hex.Dump(encoded))
	}
}
