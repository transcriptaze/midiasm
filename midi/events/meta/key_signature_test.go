package metaevent

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"testing"
)

func TestParseCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		MetaEvent{
			"KeySignature",
			events.Event{0xff, []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x00}},
			0x59,
		}, 0, 0, "C major"}

	ctx := context.NewContext()
	e := events.Event{
		Status: types.Status(0xff),
		Bytes:  []byte{0x00, 0xff},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x59, 0x02, 0x00, 0x00}))

	event, err := Parse(e, r, ctx)
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
		MetaEvent{
			"KeySignature",
			events.Event{0xff, []byte{0x00, 0xff, 0x59, 0x02, 0xfd, 0x01}},
			0x59,
		}, -3, 1, "C minor"}

	ctx := context.Context{}

	e := events.Event{
		Status: types.Status(0xff),
		Bytes:  []byte{0x00, 0xff},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x59, 0x02, 0xfd, 0x01}))

	event, err := Parse(e, r, &ctx)
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
