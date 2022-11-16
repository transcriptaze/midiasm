package disassemble

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/transcriptaze/midiasm/encoding/midi"
)

//go:embed test-files/reference.txt
var reference string

//go:embed test-files/reference-template.txt
var templated string

// //go:embed test-files/reference.mid
// var midi []byte

func TestDisassembleSMF(t *testing.T) {
	expected := reference

	midi := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x29,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27,
		0x00, 0xff, 0x2f, 0x00,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0xf1,
		0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
		0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
		0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
		0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
		0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f,
		0x00, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61,
		0x00, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73,
		0x00, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c,
		0x00, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65,
		0x00, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67,
		0x00, 0xff, 0x20, 0x01, 0x0d,
		0x00, 0xff, 0x21, 0x01, 0x70,
		0x00, 0xff, 0x59, 0x02, 0x00, 0x01,
		0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
		0x00, 0xb0, 0x00, 0x05,
		0x00, 0xb0, 0x20, 0x21,
		0x00, 0xc0, 0x19,
		0x00, 0xb0, 0x65, 0x00,
		0x00, 0xa0, 0x64,
		0x00, 0xd0, 0x07,
		0x00, 0x90, 0x30, 0x48,
		0x00, 0x92, 0x31, 0x48,
		0x00, 0x30, 0x64,
		0x81, 0x70, 0xe0, 0x00, 0x08,
		0x83, 0x60, 0x80, 0x30, 0x40,
		0x00, 0xf0, 0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7,
		0x00, 0xf0, 0x03, 0x43, 0x12, 0x00,
		0x81, 0x48, 0xF7, 0x06, 0x43, 0x12, 0x00, 0x43, 0x12, 0x00,
		0x64, 0xF7, 0x04, 0x43, 0x12, 0x00, 0xF7,
		0x00, 0xF7, 0x02, 0xF3, 0x01,
		0x00, 0xff, 0x2f, 0x00,
	}

	smf, err := midifile.NewDecoder().Decode(bytes.NewReader(midi))
	if err != nil {
		t.Fatalf("Error unmarshaling SMF: %v", err)
	}

	var s strings.Builder

	disassemble, err := NewDisassemble()
	if err != nil {
		t.Fatalf("Unexpected error initialising 'print' operation (%v)", err)
	}

	disassemble.Print(smf, "document", &s)

	if s.String() != expected {
		l, ls, p, q := diff(expected, s.String())
		t.Errorf("Output does not match expected:\n%s\n>> line %d:\n>> %s\n--------\n   %s\n   %s\n--------\n", s.String(), l, ls, p, q)
		diff(expected, s.String())
	}
}

func TestDisassembleWithLoadedTemplate(t *testing.T) {
	expected := templated

	midi := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x29,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27,
		0x00, 0xff, 0x2f, 0x00,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0xea,
		0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
		0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
		0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
		0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
		0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f,
		0x00, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61,
		0x00, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73,
		0x00, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c,
		0x00, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65,
		0x00, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67,
		0x00, 0xff, 0x20, 0x01, 0x0d,
		0x00, 0xff, 0x21, 0x01, 0x70,
		0x00, 0xff, 0x59, 0x02, 0x00, 0x01,
		0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
		0x00, 0xb0, 0x00, 0x05,
		0x00, 0xb0, 0x20, 0x21,
		0x00, 0xc0, 0x19,
		0x00, 0xb0, 0x65, 0x00,
		0x00, 0xa0, 0x64,
		0x00, 0xd0, 0x07,
		0x00, 0x90, 0x30, 0x48,
		0x81, 0x70, 0xe0, 0x00, 0x08,
		0x83, 0x60, 0x80, 0x30, 0x40,
		0x00, 0xf0, 0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7,
		0x00, 0xf0, 0x03, 0x43, 0x12, 0x00,
		0x81, 0x48, 0xF7, 0x06, 0x43, 0x12, 0x00, 0x43, 0x12, 0x00,
		0x64, 0xF7, 0x04, 0x43, 0x12, 0x00, 0xF7,
		0x00, 0xF7, 0x02, 0xF3, 0x01,
		0x00, 0xff, 0x2f, 0x00,
	}

	template := `{
  "templates": {
    "trackname": "{{.Type}} {{pad 22 .Tag}} >>> {{.Name}}"
  }
}`

	smf, err := midifile.NewDecoder().Decode(bytes.NewReader(midi))
	if err != nil {
		t.Fatalf("Error unmarshaling SMF: %v", err)
	}

	var s strings.Builder

	disassemble, err := NewDisassemble()
	if err != nil {
		t.Fatalf("Unexpected error initialising 'print' operation (%v)", err)
	}

	r := strings.NewReader(template)
	err = disassemble.LoadTemplates(r)
	if err != nil {
		t.Fatalf("Unexpected error loading 'print' templates (%v)", err)
	}

	disassemble.Print(smf, "document", &s)

	if s.String() != expected {
		l, ls, p, q := diff(expected, s.String())
		t.Errorf("Output does not match expected:\n%s\n>> line %d:\n>> %s\n--------\n   %v\n   %v\n--------\n", s.String(), l, ls, p, q)
		diff(expected, s.String())
	}
}

func TestEllipsize(t *testing.T) {
	list := map[int]string{
		0:  "",
		1:  `…`,
		2:  `1…`,
		7:  `12 34 …`,
		8:  `12 34 56`,
		16: `12 34 56`,
	}

	for l, expected := range list {
		if s := ellipsize("12 34 56", l); s != expected {
			t.Errorf("Ellipsized output incorrect - expected: <%s>, got: <%s>", expected, s)
		}
	}

	list = map[int]string{
		7: `1… 34 …`,
		8: `1… 34 56`,
	}

	for l, expected := range list {
		if s := ellipsize("1… 34 56", l); s != expected {
			t.Errorf("Ellipsized output incorrect - expected: <%s>, got: <%s>", expected, s)
		}
	}
}

func diff(p, q string) (int, string, string, string) {
	line := 0
	s := strings.Split(p, "\n")
	t := strings.Split(q, "\n")

	for ix := 0; ix < len(s) && ix < len(t); ix++ {
		line++
		if s[ix] != t[ix] {
			u := []rune(s[ix])
			v := []rune(t[ix])
			for jx := 0; jx < len(u) && jx < len(v); jx++ {
				if u[jx] != v[jx] {
					break
				}
				u[jx] = '.'
				v[jx] = '.'
			}

			return line, s[ix], string(u), string(v)
		}
	}

	return line, "?", "?", "?"
}
