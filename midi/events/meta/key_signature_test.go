package metaevent

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   0,
			delta:  0,
			Tag:    "KeySignature",
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x00},
		},
		Accidentals: 0,
		KeyType:     types.Major,
		Key:         "C major",
	}

	ctx := context.NewContext()
	r := bufio.NewReader(bytes.NewReader([]byte{0xff, 0x59, 0x02, 0x00, 0x00}))

	event, err := Parse(ctx, r, 0, 0)
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
		event: event{
			tick:   0,
			delta:  0,
			Tag:    "KeySignature",
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x00, 0xff, 0x59, 0x02, 0xfd, 0x01},
		},
		Accidentals: -3,
		KeyType:     types.Minor,
		Key:         "C minor",
	}

	ctx := context.Context{}
	r := bufio.NewReader(bytes.NewReader([]byte{0xff, 0x59, 0x02, 0xfd, 0x01}))

	event, err := Parse(&ctx, r, 0, 0)
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

func TestKeySignatureMarshalBinary(t *testing.T) {
	evt := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			Tag:    "KeySignature",
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{},
		},
		Accidentals: -3,
		KeyType:     types.Minor,
		Key:         "C minor",
	}

	expected := []byte{0xff, 0x59, 0x02, 0xfd, 0x01}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding KeySignature (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded KeySignature\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestKeySignatureUnmarshalText(t *testing.T) {
	text := "      00 FF 59 02 00 01                     tick:0          delta:480        59 KeySignature           B minor"
	expected := KeySignature{
		event: event{
			tick:   0,
			delta:  480,
			Tag:    "KeySignature",
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{},
		},
		Accidentals: 2,
		KeyType:     1,
		Key:         "B minor",
	}

	evt := KeySignature{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling KeySignature (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled KeySignature\n   expected:%#v\n   got:     %#v", expected, evt)
	}

}
