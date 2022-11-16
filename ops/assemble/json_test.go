package assemble

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"reflect"
	"testing"
)

//go:embed test-files/reference.json
var referenceJ []byte

//go:embed test-files/json.mid
var smfJ []byte

func TestJSONReference(t *testing.T) {
	assembler := JSONAssembler{}

	encoded, err := assembler.Assemble(bytes.NewBuffer(referenceJ))
	if err != nil {
		t.Fatalf("error assembling JSON file (%v)", err)
	}

	if !reflect.DeepEqual(encoded, smfJ) {
		t.Errorf("incorrectly assembled JSON file\nexpected:\n%+v\ngot:\n%+v", hex.Dump(smfJ), hex.Dump(encoded))
	}
}
