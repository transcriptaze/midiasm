package lib

import (
	"reflect"
	"testing"
)

func TestManufacturerMarshalJSON(t *testing.T) {
	manufacturer := Manufacturer{
		ID:     []byte{0x00, 0x00, 0x3b},
		Region: "American",
		Name:   "Mark Of The Unicorn (MOTU)",
	}

	expected := `{"id":[0,0,59],"region":"American","name":"Mark Of The Unicorn (MOTU)"}`

	if bytes, err := manufacturer.MarshalJSON(); err != nil {
		t.Errorf("Error marshalling Manufacturer (%v)", err)
	} else if string(bytes) != expected {
		t.Errorf("Incorrectly marshalled Manufacturer - expected:%v, got:%v", expected, string(bytes))
	}
}

func TestManufacturerUnmarshalJSON(t *testing.T) {
	bytes := `{"id":[0,0,59],"region":"American","name":"Mark Of The Unicorn (MOTU)"}`
	expected := Manufacturer{
		ID:     []byte{0x00, 0x00, 0x3b},
		Region: "American",
		Name:   "Mark Of The Unicorn (MOTU)",
	}

	var m Manufacturer

	if err := m.UnmarshalJSON([]byte(bytes)); err != nil {
		t.Errorf("Error unmarshalling Manufacturer (%v)", err)
	} else if !reflect.DeepEqual(m, expected) {
		t.Errorf("Incorrectly marshalled Manufacturer - expected:%#v, got:%#v", expected, m)
	}
}
