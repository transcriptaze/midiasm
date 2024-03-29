package sysex

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParseSysExSingleMessage(t *testing.T) {
	event, err := Parse(0, false, []byte{0x83, 0x60, 0xf0, 0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}...)
	if err != nil {
		t.Fatalf("Unexpected SysEx single message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx single message parse error - returned %v", event)
	}

	message, ok := event.(SysExMessage)
	if !ok {
		t.Fatalf("SysEx single message parse error - returned %T", event)
	}

	manufacturer := lib.Manufacturer{
		ID:     []byte{0x7e},
		Region: "Special Purpose",
		Name:   "Non-RealTime Extensions",
	}
	if !reflect.DeepEqual(message.Manufacturer, manufacturer) {
		t.Errorf("Invalid SysEx single message manufacturer - expected:%v, got: %v", manufacturer, message.Manufacturer)
	}

	data := lib.Hex{0x00, 0x09, 0x01}
	if !reflect.DeepEqual(message.Data, data) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", data, message.Data)
	}

	if !message.Single {
		t.Errorf("SysEx single message 'Single' flag missing - expected:%v, got: %v", true, message.Single)
	}
}

func TestParseSysExMessage(t *testing.T) {
	event, err := Parse(0, false, []byte{0x83, 0x60, 0xf0, 0x05, 0x7e, 0x00, 0x09, 0x01, 0x43}...)
	if err != nil {
		t.Fatalf("Unexpected SysEx message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx message parse error - returned %v", event)
	}

	message, ok := event.(SysExMessage)
	if !ok {
		t.Fatalf("SysEx message parse error - returned %T", event)
	}

	manufacturer := lib.Manufacturer{
		ID:     []byte{0x7e},
		Region: "Special Purpose",
		Name:   "Non-RealTime Extensions",
	}
	if !reflect.DeepEqual(message.Manufacturer, manufacturer) {
		t.Errorf("Invalid SysEx message manufacturer - expected:%v, got: %v", manufacturer, message.Manufacturer)
	}

	data := lib.Hex([]byte{0x00, 0x09, 0x01, 0x43})
	if !reflect.DeepEqual(message.Data, data) {
		t.Errorf("Invalid SysEx message data - expected:%v, got: %v", data, message.Data)
	}

	// if !ctx.Casio {
	// 	t.Errorf("context Casio flag should be set")
	// }
}

func TestSysExMessageMarshalBinary(t *testing.T) {
	evt := SysExMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data: lib.Hex{0x00, 0x09, 0x01},
	}

	expected := []byte{0xf0, 0x04, 0x7e, 0x00, 0x09, 0x01}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SysExMessage (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SysExMessage\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSysExMessageUnmarshalBinary(t *testing.T) {
	expected := SysExMessage{
		event: event{
			delta:  480,
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
			bytes:  []byte{0x83, 0x60, 0xf0, 0x04, 0x7e, 0x00, 0x09, 0x01},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data: lib.Hex{0x00, 0x09, 0x01},
	}

	bytes := []byte{0x83, 0x60, 0xf0, 0x04, 0x7e, 0x00, 0x09, 0x01}

	e := SysExMessage{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error unencoding %v (%v)", lib.TagSysExMessage, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSysExMessage, expected, e)
	}
}

func TestSysExSingleMessageMarshalBinary(t *testing.T) {
	evt := SysExMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data:   lib.Hex{0x00, 0x09, 0x01},
		Single: true,
	}

	expected := []byte{0xf0, 0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SysExMessage (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SysExMessage\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSysExMessageUnmarshalText(t *testing.T) {
	text := "      00 F0 05 7E 00 09 01 F7               tick:720        delta:480         F0 SysExMessage           Non-RealTime Extensions, 00 09 01"
	expected := SysExMessage{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
			bytes:  []byte{},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data: lib.Hex{0x00, 0x09, 0x01},
	}

	evt := SysExMessage{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SysExMessage (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SysExMessage\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestSysExMessageMarshalJSON(t *testing.T) {
	e := SysExMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data:   lib.Hex{0x00, 0x09, 0x01},
		Single: true,
	}

	expected := `{"tag":"SysExMessage","delta":480,"status":240,"manufacturer":{"id":[126],"region":"Special Purpose","name":"Non-RealTime Extensions"},"data":[0,9,1],"single":true}`

	testMarshalJSON(t, lib.TagSysExMessage, e, expected)
}

func TestSysExMessageNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagSysExMessage
	text := `{"tag":"SysExMessage","delta":480,"status":240,"manufacturer":{"id":[126],"region":"Special Purpose","name":"Non-RealTime Extensions"},"data":[0,9,1],"single":true}`
	expected := SysExMessage{
		event: event{
			tick:   0,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x7e},
			Region: "Special Purpose",
			Name:   "Non-RealTime Extensions",
		},
		Data:   lib.Hex{0x00, 0x09, 0x01},
		Single: true,
	}

	e := SysExMessage{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
