package lib

import (
	"fmt"
)

type Note struct {
	Ord  int
	Name string
}

func (n Note) String() string {
	return n.Name
}

func LookupNote(s string) (Note, error) {
	for _, note := range Notes {
		if note.Name == s {
			return note, nil
		}
	}

	return Note{}, fmt.Errorf("invalid note (%v)", s)
}

var Notes = []Note{
	C,
	C_SHARP,
	D_FLAT,
	D,
	D_SHARP,
	E_FLAT,
	E,
	E_SHARP,
	F_FLAT,
	F,
	F_SHARP,
	G_FLAT,
	G,
	G_SHARP,
	A_FLAT,
	A,
	A_SHARP,
	B_FLAT,
	B,
	B_SHARP,
	C_FLAT,
}

var C = Note{
	Ord:  0,
	Name: `C`,
}

var C_SHARP = Note{
	Ord:  1,
	Name: `C♯`,
}

var D_FLAT = Note{
	Ord:  1,
	Name: `D♭`,
}

var D = Note{
	Ord:  2,
	Name: `D`,
}

var D_SHARP = Note{
	Ord:  3,
	Name: `D♯`,
}

var E_FLAT = Note{
	Ord:  3,
	Name: `E♭`,
}

var E = Note{
	Ord:  4,
	Name: `E`,
}

var E_SHARP = Note{
	Ord:  5,
	Name: `E♯`,
}

var F_FLAT = Note{
	Ord:  4,
	Name: `F♭`,
}

var F = Note{
	Ord:  5,
	Name: `F`,
}

var F_SHARP = Note{
	Ord:  6,
	Name: `F♯`,
}

var G_FLAT = Note{
	Ord:  6,
	Name: `G♭`,
}

var G = Note{
	Ord:  7,
	Name: `G`,
}

var G_SHARP = Note{
	Ord:  8,
	Name: `G♯`,
}

var A_FLAT = Note{
	Ord:  8,
	Name: `A♭`,
}

var A = Note{
	Ord:  9,
	Name: `A`,
}

var A_SHARP = Note{
	Ord:  10,
	Name: `A♯`,
}

var B_FLAT = Note{
	Ord:  10,
	Name: `B♭`,
}

var B = Note{
	Ord:  11,
	Name: `B`,
}

var B_SHARP = Note{
	Ord:  0,
	Name: `B♯`,
}

var C_FLAT = Note{
	Ord:  11,
	Name: `C♭`,
}
