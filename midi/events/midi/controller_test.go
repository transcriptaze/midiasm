package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestController(t *testing.T) {
	expected := Controller{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xb7, 0x54, 0x1d},

			tag:     types.TagController,
			Status:  0xb7,
			Channel: 7,
		},
		Controller: types.Controller{84, "Portamento Control"},
		Value:      29,
	}

	ctx := context.NewContext()
	r := IO.BytesReader([]byte{0x54, 0x1d})

	event, err := Parse(2400, 480, r, 0xb7, ctx)
	if err != nil {
		t.Fatalf("Unexpected Controller event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected Controller event parse error - returned %v", event)
	}

	event, ok := event.(*Controller)
	if !ok {
		t.Fatalf("Controller event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid Controller event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestProgramBank(t *testing.T) {
	ctx := context.NewContext()
	r := IO.BytesReader([]byte{0x00, 0x05, 0x20, 0x21})

	if _, err := Parse(0, 0, r, 0xb3, ctx); err != nil {
		t.Fatalf("Unexpected MIDI event parse error: %v", err)
	}

	if _, err := Parse(0, 0, r, 0xb3, ctx); err != nil {
		t.Fatalf("Unexpected MIDI event parse error: %v", err)
	}

	if ctx.ProgramBank[3] != 673 {
		t.Errorf("Invalid ProgramBank in context\n  expected:%v\n  got:     %#v", 673, ctx.ProgramBank[3])
	}
}

func TestControllerMarshalBinary(t *testing.T) {
	evt := Controller{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xb7, 0x54, 0x1d},
			tag:   types.TagController,

			Status:  0xb7,
			Channel: 7,
		},
		Controller: types.Controller{84, "Portamento Control"},
		Value:      29,
	}

	expected := []byte{0xb7, 0x54, 0x1d}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Controller (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Controller\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestControllerUnmarshalText(t *testing.T) {
	text := "      00 B0 65 09                           tick:0          delta:480        B7 Controller             channel:7  101/Registered Parameter Number (MSB), value:9"
	expected := Controller{
		event: event{
			tick:    0,
			delta:   480,
			tag:     types.TagController,
			Status:  0xb7,
			Channel: 7,
			bytes:   []byte{},
		},
		Controller: types.Controller{101, "Registered Parameter Number (MSB)"},
		Value:      0x09,
	}

	evt := Controller{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Controller (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Controller\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
