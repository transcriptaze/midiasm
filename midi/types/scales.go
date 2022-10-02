package types

import (
	"fmt"
)

type Scale struct {
	Name        string
	Accidentals int8
	Type        KeyType
	Notes       []Note
}

type KeyType uint8

const (
	Major KeyType = 0
	Minor KeyType = 1
)

type Note struct {
	Name string
	Ord  int
}

func (n Note) String() string {
	return n.Name
}

func (s Scale) Transpose(steps int) Scale {
	notes := make([]Note, len(s.Notes))

	// fmt.Printf(">>>>>>>>> SCALE: %+v\n", s.Notes)

	for ix, n := range s.Notes {
		notes[ix] = transpose(n, steps)
	}

	// fmt.Printf(">>>>>>>>> NOTES: %+v\n", notes)
loop:
	for _, scale := range MAJOR_SCALES {
		for ix, note := range notes {
			m := scale.Notes[ix]

			if note.Ord == m.Ord {
				continue
			}

			continue loop
		}

		return scale
	}

	return s
}

func transpose(note Note, steps int) Note {
	for ix, n := range SCALE {
		if n.Ord == note.Ord {
			ix = (ix + steps)
			for ix < 0 {
				ix += 12
			}

			ix %= len(SCALE)

			return SCALE[ix]
		}
	}

	panic(fmt.Sprintf("invalid notes %q", note))
}

var SCALE = []Note{
	C, C_SHARP, D, D_SHARP, E, F, F_SHARP, G, G_SHARP, A, A_SHARP, B,
}

// {C, D_FLAT, D, E_FLAT, E, F, G_FLAT, G, A_FLAT, A, B_FLAT, B},

func MajorScale(accidentals int8) (Scale, bool) {
	for _, scale := range MAJOR_SCALES {
		if scale.Accidentals == accidentals {
			return scale, true
		}
	}

	return C_MAJOR, false
}

func MinorScale(accidentals int8) (Scale, bool) {
	for _, scale := range MINOR_SCALES {
		if scale.Accidentals == accidentals {
			return scale, true
		}
	}

	return A_MINOR, false
}

var MAJOR_SCALES = []Scale{
	C_MAJOR,
	G_MAJOR,
	D_MAJOR,
	A_MAJOR,
	E_MAJOR,
	B_MAJOR,
	F_SHARP_MAJOR,
	C_SHARP_MAJOR,
	F_MAJOR,
	B_FLAT_MAJOR,
	E_FLAT_MAJOR,
	A_FLAT_MAJOR,
	D_FLAT_MAJOR,
	G_FLAT_MAJOR,
	C_FLAT_MAJOR,
}

var MINOR_SCALES = []Scale{
	A_SHARP_MINOR,
	D_SHARP_MINOR,
	G_SHARP_MINOR,
	C_SHARP_MINOR,
	F_SHARP_MINOR,
	B_MINOR,
	E_MINOR,
	A_MINOR,
	D_MINOR,
	G_MINOR,
	C_MINOR,
	F_MINOR,
	B_FLAT_MINOR,
	E_FLAT_MINOR,
	A_FLAT_MINOR,
}

var C = Note{Name: `C`, Ord: 0}
var C_SHARP = Note{Name: `C♯`, Ord: 1}
var D_FLAT = Note{Name: `D♭`, Ord: 1}
var D = Note{Name: `D`, Ord: 2}
var D_SHARP = Note{Name: `D♯`, Ord: 3}
var E_FLAT = Note{Name: `E♭`, Ord: 3}
var E = Note{Name: `E`, Ord: 4}
var E_SHARP = Note{Name: `E♯`, Ord: 5}
var F_FLAT = Note{Name: `F♭`, Ord: 4}
var F = Note{Name: `F`, Ord: 5}
var F_SHARP = Note{Name: `F♯`, Ord: 6}
var G_FLAT = Note{Name: `G♭`, Ord: 6}
var G = Note{Name: `G`, Ord: 7}
var G_SHARP = Note{Name: `G♯`, Ord: 8}
var A_FLAT = Note{Name: `A♭`, Ord: 8}
var A = Note{Name: `A`, Ord: 9}
var A_SHARP = Note{Name: `A♯`, Ord: 10}
var B_FLAT = Note{Name: `B♭`, Ord: 10}
var B = Note{Name: `B`, Ord: 11}
var B_SHARP = Note{Name: `B♯`, Ord: 0}
var C_FLAT = Note{Name: `C♭`, Ord: 11}

var C_MAJOR = Scale{
	Name:        `C major`,
	Accidentals: 0,
	Type:        Major,
	Notes:       []Note{C, D, E, F, G, A, B},
}

var G_MAJOR = Scale{
	Name:        `G major`,
	Accidentals: 1,
	Type:        Major,
	Notes:       []Note{G, A, B, C, D, E, F_SHARP},
}

var D_MAJOR = Scale{
	Name:        `D major`,
	Accidentals: 2,
	Type:        Major,
	Notes:       []Note{D, E, F_SHARP, G, A, B, C_SHARP},
}

var A_MAJOR = Scale{
	Name:        `A major`,
	Accidentals: 3,
	Type:        Major,
	Notes:       []Note{A, B, C_SHARP, D, E, F_SHARP, G_SHARP},
}

var E_MAJOR = Scale{
	Name:        `E major`,
	Accidentals: 4,
	Type:        Major,
	Notes:       []Note{E, F_SHARP, G_SHARP, A, B, C_SHARP, D_SHARP},
}

var B_MAJOR = Scale{
	Name:        `B major`,
	Accidentals: 5,
	Type:        Major,
	Notes:       []Note{B, C_SHARP, D_SHARP, E, F_SHARP, G_SHARP, A_SHARP},
}

var F_SHARP_MAJOR = Scale{
	Name:        `F♯ major`,
	Accidentals: 6,
	Type:        Major,
	Notes:       []Note{F_SHARP, G_SHARP, A_SHARP, B, C_SHARP, D_SHARP, E_SHARP},
}

var C_SHARP_MAJOR = Scale{
	Name:        `C♯ major`,
	Accidentals: 7,
	Type:        Major,
	Notes:       []Note{C_SHARP, D_SHARP, E_SHARP, F_SHARP, G_SHARP, A_SHARP, B_SHARP},
}

var F_MAJOR = Scale{
	Name:        `F major`,
	Accidentals: -1,
	Type:        Major,
	Notes:       []Note{F, G, A, B_FLAT, C, D, E},
}

var B_FLAT_MAJOR = Scale{
	Name:        `B♭ major`,
	Accidentals: -2,
	Type:        Major,
	Notes:       []Note{B_FLAT, C, D, E_FLAT, F, G, A},
}

var E_FLAT_MAJOR = Scale{
	Name:        `E♭ major`,
	Accidentals: -3,
	Type:        Major,
	Notes:       []Note{E_FLAT, F, G, A_FLAT, B_FLAT, C, D},
}

var A_FLAT_MAJOR = Scale{
	Name:        `A♭ major`,
	Accidentals: -4,
	Type:        Major,
	Notes:       []Note{A_FLAT, B_FLAT, C, D_FLAT, E_FLAT, F, G},
}

var D_FLAT_MAJOR = Scale{
	Name:        `D♭ major`,
	Accidentals: -5,
	Type:        Major,
	Notes:       []Note{C, D_FLAT, E_FLAT, F, G_FLAT, A_FLAT, B_FLAT},
}

var G_FLAT_MAJOR = Scale{
	Name:        `G♭ major`,
	Accidentals: -6,
	Type:        Major,
	Notes:       []Note{C_FLAT, D_FLAT, E_FLAT, F, G_FLAT, A_FLAT, B_FLAT},
}

var C_FLAT_MAJOR = Scale{
	Name:        `C♭ major`,
	Accidentals: -7,
	Type:        Major,
	Notes:       []Note{C_FLAT, D_FLAT, E_FLAT, F_FLAT, G_FLAT, A_FLAT, B_FLAT},
}

var A_SHARP_MINOR = Scale{
	Name:        `A♯ minor`,
	Accidentals: 7,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D_SHARP, E_SHARP, F_SHARP, G_SHARP, A_SHARP, B_SHARP},
}

var D_SHARP_MINOR = Scale{
	Name:        `D♯ minor`,
	Accidentals: 6,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D_SHARP, E_SHARP, F_SHARP, G_SHARP, A_SHARP, B},
}

var G_SHARP_MINOR = Scale{
	Name:        `G♯ minor`,
	Accidentals: 5,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D_SHARP, E, F_SHARP, G_SHARP, A_SHARP, B},
}

var C_SHARP_MINOR = Scale{
	Name:        `C♯ minor`,
	Accidentals: 4,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D_SHARP, E, F_SHARP, G_SHARP, A, B},
}

var F_SHARP_MINOR = Scale{
	Name:        `F♯ minor`,
	Accidentals: 3,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D, E, F_SHARP, G_SHARP, A, B},
}

var B_MINOR = Scale{
	Name:        `B minor`,
	Accidentals: 2,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D, E, F_SHARP, G, A, B},
}

var E_MINOR = Scale{
	Name:        `E minor`,
	Accidentals: 1,
	Type:        Minor,
	Notes:       []Note{C, D, E, F_SHARP, G, A, B},
}

var A_MINOR = Scale{
	Name:        `A minor`,
	Accidentals: 0,
	Type:        Minor,
	Notes:       []Note{C, D, E, F, G, A, B},
}

var D_MINOR = Scale{
	Name:        `D minor`,
	Accidentals: -1,
	Type:        Minor,
	Notes:       []Note{C, D, E, F, G, A, B_FLAT},
}

var G_MINOR = Scale{
	Name:        `G minor`,
	Accidentals: -2,
	Type:        Minor,
	Notes:       []Note{C, D, E_FLAT, F, G, A, B_FLAT},
}

var C_MINOR = Scale{
	Name:        `C minor`,
	Accidentals: -3,
	Type:        Minor,
	Notes:       []Note{C, D, E_FLAT, F, G, A_FLAT, B_FLAT},
}

var F_MINOR = Scale{
	Name:        `F minor`,
	Accidentals: -4,
	Type:        Minor,
	Notes:       []Note{C, D_FLAT, E_FLAT, F, G, A_FLAT, B_FLAT},
}

var B_FLAT_MINOR = Scale{
	Name:        `B♭ minor`,
	Accidentals: -5,
	Type:        Minor,
	Notes:       []Note{C, D_FLAT, E_FLAT, F, G_FLAT, A_FLAT, B_FLAT},
}

var E_FLAT_MINOR = Scale{
	Name:        `E♭ minor`,
	Accidentals: -6,
	Type:        Minor,
	Notes:       []Note{C_FLAT, D_FLAT, E_FLAT, F, G_FLAT, A_FLAT, B_FLAT},
}

var A_FLAT_MINOR = Scale{
	Name:        `A♭ minor`,
	Accidentals: -7,
	Type:        Minor,
	Notes:       []Note{C_FLAT, D_FLAT, E_FLAT, F_FLAT, G_FLAT, A_FLAT, B_FLAT},
}
