package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x00},
		},
		Accidentals: 0,
		KeyType:     types.Major,
		Key:         "C major",
	}

	ctx := context.NewContext()

	event, err := Parse(ctx, 2400, 480, 0xff, 0x59, []byte{0x00, 0x00})
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if ks, ok := event.(*KeySignature); !ok {
		t.Fatalf("KeySignature event parse error - returned %T", event)
	} else if !reflect.DeepEqual(ks, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &expected, ks)
	}

	if !reflect.DeepEqual(ctx.Scale(), context.Sharps) {
		t.Errorf("Context scale not set to 'sharps':%v", ctx)
	}
}

func TestParseCMinorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x00, 0xff, 0x59, 0x02, 0xfd, 0x01},
		},
		Accidentals: -3,
		KeyType:     types.Minor,
		Key:         "C minor",
	}

	ctx := context.NewContext()

	event, err := Parse(ctx, 2400, 480, 0xff, 0x59, []byte{0xfd, 0x01})
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if ks, ok := event.(*KeySignature); !ok {
		t.Fatalf("KeySignature event parse error - returned %T", event)
	} else if !reflect.DeepEqual(ks, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &ks, event)
	}

	if !reflect.DeepEqual(ctx.Scale(), context.Flats) {
		t.Errorf("Context scale not set to 'flats':%v", ctx)
	}
}
