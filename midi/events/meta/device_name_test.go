package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalDeviceName(t *testing.T) {
	expected := DeviceName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   0x09,
			bytes:  []byte{0x00, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67},
		},
		Name: "TheThing",
	}

	e := DeviceName{}

	err := e.unmarshal(2400, 480, 0xff, []byte("TheThing"), []byte{0x00, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67}...)
	if err != nil {
		t.Fatalf("error unmarshalling DeviceName (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect DeviceName\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestDeviceNameMarshalBinary(t *testing.T) {
	evt := DeviceName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   0x09,
			bytes:  []byte{},
		},
		Name: "TheThing",
	}

	expected := []byte{0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding DeviceName (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded DeviceName\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestDeviceNameUnmarshalBinary(t *testing.T) {
	expected := DeviceName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   lib.TypeDeviceName,
			bytes:  []byte{0x83, 0x60, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67},
		},
		Name: "TheThing",
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67}

	e := DeviceName{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagDeviceName, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagDeviceName, expected, e)
	}
}

func TestDeviceNameUnmarshalText(t *testing.T) {
	text := "      00 FF 09 08 54 68 65 54 68 69 6E 67   tick:0          delta:480        09 DeviceName             TheThing"
	expected := DeviceName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   0x09,
			bytes:  []byte{},
		},
		Name: "TheThing",
	}

	evt := DeviceName{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling DeviceName (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled DeviceName\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestDeviceNameMarshalJSON(t *testing.T) {
	e := DeviceName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   lib.TypeDeviceName,
			bytes:  []byte{},
		},
		Name: "TheThing",
	}

	expected := `{"tag":"DeviceName","delta":480,"status":255,"type":9,"name":"TheThing"}`

	testMarshalJSON(t, lib.TagDeviceName, e, expected)
}

func TestDeviceNameNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagDeviceName
	text := `{"tag":"DeviceName","delta":480,"status":255,"type":9,"name":"TheThing"}`
	expected := DeviceName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagDeviceName,
			Status: 0xff,
			Type:   lib.TypeDeviceName,
			bytes:  []byte{},
		},
		Name: "TheThing",
	}

	evt := DeviceName{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
