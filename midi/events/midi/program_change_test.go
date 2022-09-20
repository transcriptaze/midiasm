package midievent

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
)

func TestProgramChange(t *testing.T) {
	expected := ProgramChange{
		Tag:     "ProgramChange",
		Status:  0xc7,
		Channel: 7,
		Bank:    673,
		Program: 13,
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x0d}))

	ctx := context.NewContext()
	ctx.ProgramBank[7] = 673

	event, err := Parse(r, 0xc7, ctx)
	if err != nil {
		t.Fatalf("Unexpected ProgramChange event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected ProgramChange event parse error - returned %v", event)
	}

	event, ok := event.(*ProgramChange)
	if !ok {
		t.Fatalf("ProgramChange event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid ProgramChange event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}
