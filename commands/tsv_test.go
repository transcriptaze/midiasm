package commands

import (
	"bytes"
	_ "embed"
	"reflect"
	"testing"
)

//go:embed test-files/reference.mid
var _SMF []byte

//go:embed test-files/reference.tsv
var _TSV []byte

func TestTSV(t *testing.T) {
	var v = tsv{}
	var b bytes.Buffer

	if smf, err := v.decode(bytes.NewBuffer(_SMF)); err != nil {
		t.Fatalf("%v", err)
	} else if header, records, err := v.export(smf); err != nil {
		t.Fatalf("%v", err)
	} else if err := writeTSV(header, records, '\t', &b); err != nil {
		t.Fatalf("%v", err)
	} else if !reflect.DeepEqual(b.Bytes(), _TSV) {
		t.Errorf("Incorrectly exported to TSV file")
	}
}
