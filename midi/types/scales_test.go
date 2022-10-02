package types

import (
	"reflect"
	"testing"
)

func TestTransposeCMajorScale(t *testing.T) {
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

	for octave := -2; octave <= 2; octave++ {
		for _, v := range tests {
			transposed := v.scale.Transpose(12*octave + v.steps)

			if !reflect.DeepEqual(transposed, v.expected) {
				t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", v.scale.Name, v.expected, transposed)
			}
		}
	}
}

func TestTransposeAMinorScale(t *testing.T) {
	tests := []struct {
		scale    Scale
		steps    int
		expected Scale
	}{
		{A_MINOR, 1, A_SHARP_MINOR},
		{A_MINOR, 2, B_MINOR},
		{A_MINOR, 3, C_MINOR},
		{A_MINOR, 4, C_SHARP_MINOR},
		{A_MINOR, 5, D_MINOR},
		{A_MINOR, 6, D_SHARP_MINOR},
		{A_MINOR, 7, E_MINOR},
		{A_MINOR, 8, F_MINOR},
		{A_MINOR, 9, F_SHARP_MINOR},
		{A_MINOR, 10, G_MINOR},
		{A_MINOR, 11, G_SHARP_MINOR},
		{A_MINOR, 12, A_MINOR},

		{A_MINOR, -1, G_SHARP_MINOR},
		{A_MINOR, -2, G_MINOR},
		{A_MINOR, -3, F_SHARP_MINOR},
		{A_MINOR, -4, F_MINOR},
		{A_MINOR, -5, E_MINOR},
		{A_MINOR, -6, D_SHARP_MINOR},
		{A_MINOR, -7, D_MINOR},
		{A_MINOR, -8, C_SHARP_MINOR},
		{A_MINOR, -9, C_MINOR},
		{A_MINOR, -10, B_MINOR},
		{A_MINOR, -11, A_SHARP_MINOR},
		{A_MINOR, -12, A_MINOR},
	}

	for octave := -2; octave <= 2; octave++ {
		for _, v := range tests {
			transposed := v.scale.Transpose(12*octave + v.steps)

			if !reflect.DeepEqual(transposed, v.expected) {
				t.Errorf("Incorrectly transposed %v scale\n   expected:%+v\n   got:     %+v", v.scale.Name, v.expected, transposed)
			}
		}
	}
}
