package types

import (
	"reflect"
	"testing"
)

func TestScaleTranspose(t *testing.T) {
	tests := []struct {
		scale    Scale
		expected Scale
	}{
		{C_MAJOR, C_SHARP_MAJOR},
		// {G_MAJOR, C_SHARP_MAJOR},
	}

	for _, v := range tests {
		scale := v.scale
		expected := v.expected
		transposed := scale.Transpose(1)

		if !reflect.DeepEqual(transposed, expected) {
			t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", scale.Name, expected, transposed)
		}
	}
}

func TestScaleTransposeByFifths(t *testing.T) {
	tests := []struct {
		scale    Scale
		expected Scale
	}{
		{C_MAJOR, F_MAJOR},
	}

	for _, v := range tests {
		scale := v.scale
		expected := v.expected
		transposed := scale.Transpose(5)

		if !reflect.DeepEqual(transposed, expected) {
			t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", scale.Name, expected, transposed)
		}
	}
}
