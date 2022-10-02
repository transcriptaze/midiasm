package types

import (
	"reflect"
	"testing"
)

func TestCMajorScaleTransposeUp(t *testing.T) {
	tests := []struct {
		scale    Scale
		steps    int
		expected Scale
	}{
		{C_MAJOR, 1, C_SHARP_MAJOR},
		{C_MAJOR, 2, D_MAJOR},
		{C_MAJOR, 3, E_FLAT_MAJOR},
		{C_MAJOR, 4, E_MAJOR},
		{C_MAJOR, 5, F_MAJOR},
		{C_MAJOR, 6, F_SHARP_MAJOR},
		{C_MAJOR, 7, G_MAJOR},
		{C_MAJOR, 8, A_FLAT_MAJOR},
		{C_MAJOR, 9, A_MAJOR},
		{C_MAJOR, 10, B_FLAT_MAJOR},
		{C_MAJOR, 11, B_MAJOR},
		{C_MAJOR, 12, C_MAJOR},
	}

	for octave := 0; octave < 8; octave++ {
		for _, v := range tests {
			transposed := v.scale.Transpose(12*octave + v.steps)

			if !reflect.DeepEqual(transposed, v.expected) {
				t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", v.scale.Name, v.expected, transposed)
			}
		}
	}
}

func TestCMajorScaleTransposeDown(t *testing.T) {
	tests := []struct {
		scale    Scale
		steps    int
		expected Scale
	}{
		{C_MAJOR, -1, B_MAJOR},
		{C_MAJOR, -2, B_FLAT_MAJOR},
		{C_MAJOR, -3, A_MAJOR},
		{C_MAJOR, -4, A_FLAT_MAJOR},
		{C_MAJOR, -5, G_MAJOR},
		{C_MAJOR, -6, F_SHARP_MAJOR},
		{C_MAJOR, -7, F_MAJOR},
		{C_MAJOR, -8, E_MAJOR},
		{C_MAJOR, -9, E_FLAT_MAJOR},
		{C_MAJOR, -10, D_MAJOR},
		{C_MAJOR, -11, C_SHARP_MAJOR},
		{C_MAJOR, -12, C_MAJOR},
	}

	for octave := 0; octave < 1; octave++ {
		for _, v := range tests {
			transposed := v.scale.Transpose(12*octave + v.steps)

			if !reflect.DeepEqual(transposed, v.expected) {
				t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", v.scale.Name, v.expected, transposed)
			}
		}
	}
}
