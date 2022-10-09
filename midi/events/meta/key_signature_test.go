package metaevent

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
)

func TestParseCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		"KeySignature",
		0xff,
		0x59, 0, 0, "C major"}

	ctx := context.NewContext()
	r := bufio.NewReader(bytes.NewReader([]byte{0x59, 0x02, 0x00, 0x00}))

	event, err := Parse(ctx, r, 0xff, 0, 0)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	event, ok := event.(*KeySignature)
	if !ok {
		t.Fatalf("KeySignature event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &expected, event)
	}

	if !reflect.DeepEqual(ctx.Scale(), context.Sharps) {
		t.Errorf("Context scale not set to 'sharps':%v", ctx)
	}
}

func TestParseCMinorKeySignature(t *testing.T) {
	expected := KeySignature{
		"KeySignature",
		0xff,
		0x59, -3, 1, "C minor"}

	ctx := context.Context{}
	r := bufio.NewReader(bytes.NewReader([]byte{0x59, 0x02, 0xfd, 0x01}))

	event, err := Parse(&ctx, r, 0xff, 0, 0)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	event, ok := event.(*KeySignature)
	if !ok {
		t.Fatalf("KeySignature event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &expected, event)
	}

	if !reflect.DeepEqual(ctx.Scale(), context.Flats) {
		t.Errorf("Context scale not set to 'flats':%v", ctx)
	}
}
