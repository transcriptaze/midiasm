package types

import (
	"fmt"
)

type ScaleID int

const (
	CMajor ScaleID = iota
	GMajor
	DMajor
	AMajor
	EMajor
	BMajor
	FSharpMajor
	CSharpMajor
	FMajor
	BFlatMajor
	EFlatMajor
	AFlatMajor
	DFlatMajor
	GFlatMajor
	CFlatMajor

	AMinor
	EMinor
	BMinor
	ASharpMinor
	DSharpMinor
	GSharpMinor
	CSharpMinor
	FSharpMinor
	DMinor
	GMinor
	CMinor
	FMinor
	BFlatMinor
	EFlatMinor
	AFlatMinor
)

type Scale struct {
	id          ScaleID
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

func (k KeyType) String() string {
	return []string{"major", "minor"}[k]
}

func (s Scale) Transpose(steps int) Scale {
	notes := make([]Note, len(s.Notes))
	for ix, n := range s.Notes {
		notes[ix] = transpose(n, steps)
	}

	var scales []Scale

	switch s.Type {
	case Major:
		scales = MAJOR_SCALES
	case Minor:
		scales = MINOR_SCALES
	default:
		return s
	}

loop:
	for _, scale := range scales {
		for ix, note := range notes {
			m := scale.Notes[ix]

			if note.Ord == m.Ord {
				continue
			}

			continue loop
		}

		// .. use enharmonic scales
		switch scale.id {
		case C_SHARP_MAJOR.id: // C♯ major/D♭ major
			return D_FLAT_MAJOR

		case G_FLAT_MAJOR.id: // /G♭ major/F♯ major
			return F_SHARP_MAJOR

		case C_FLAT_MAJOR.id: // /C♭ major/B major
			return B_MAJOR

		case A_FLAT_MINOR.id: // A♭ minor/G♯ minor
			return G_SHARP_MINOR

		case E_FLAT_MINOR.id: // E♭ minor/D♯ minor
			return D_SHARP_MINOR

		case A_SHARP_MINOR.id: // A♯ minor/B♭ minor
			return B_FLAT_MINOR
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
	A_MINOR,
	E_MINOR,
	B_MINOR,

	A_SHARP_MINOR,
	D_SHARP_MINOR,
	G_SHARP_MINOR,
	C_SHARP_MINOR,
	F_SHARP_MINOR,
	D_MINOR,
	G_MINOR,
	C_MINOR,
	F_MINOR,
	B_FLAT_MINOR,
	E_FLAT_MINOR,
	A_FLAT_MINOR,
}

// Major scales

var C_MAJOR = Scale{
	id:          CMajor,
	Name:        `C major`,
	Accidentals: 0,
	Type:        Major,
	Notes:       []Note{C, D, E, F, G, A, B},
}

var G_MAJOR = Scale{
	id:          GMajor,
	Name:        `G major`,
	Accidentals: 1,
	Type:        Major,
	Notes:       []Note{G, A, B, C, D, E, F_SHARP},
}

var D_MAJOR = Scale{
	id:          DMajor,
	Name:        `D major`,
	Accidentals: 2,
	Type:        Major,
	Notes:       []Note{D, E, F_SHARP, G, A, B, C_SHARP},
}

var A_MAJOR = Scale{
	id:          AMajor,
	Name:        `A major`,
	Accidentals: 3,
	Type:        Major,
	Notes:       []Note{A, B, C_SHARP, D, E, F_SHARP, G_SHARP},
}

var E_MAJOR = Scale{
	id:          EMajor,
	Name:        `E major`,
	Accidentals: 4,
	Type:        Major,
	Notes:       []Note{E, F_SHARP, G_SHARP, A, B, C_SHARP, D_SHARP},
}

// enharmonic to C♭ major
var B_MAJOR = Scale{
	id:          BMajor,
	Name:        `B major`,
	Accidentals: 5,
	Type:        Major,
	Notes:       []Note{B, C_SHARP, D_SHARP, E, F_SHARP, G_SHARP, A_SHARP},
}

// enharmonic to G♭ major
var F_SHARP_MAJOR = Scale{
	id:          FSharpMajor,
	Name:        `F♯ major`,
	Accidentals: 6,
	Type:        Major,
	Notes:       []Note{F_SHARP, G_SHARP, A_SHARP, B, C_SHARP, D_SHARP, E_SHARP},
}

// enharmonic to D major
var C_SHARP_MAJOR = Scale{
	id:          CSharpMajor,
	Name:        `C♯ major`,
	Accidentals: 7,
	Type:        Major,
	Notes:       []Note{C_SHARP, D_SHARP, E_SHARP, F_SHARP, G_SHARP, A_SHARP, B_SHARP},
}

var F_MAJOR = Scale{
	id:          FMajor,
	Name:        `F major`,
	Accidentals: -1,
	Type:        Major,
	Notes:       []Note{F, G, A, B_FLAT, C, D, E},
}

var B_FLAT_MAJOR = Scale{
	id:          BFlatMajor,
	Name:        `B♭ major`,
	Accidentals: -2,
	Type:        Major,
	Notes:       []Note{B_FLAT, C, D, E_FLAT, F, G, A},
}

var E_FLAT_MAJOR = Scale{
	id:          EFlatMajor,
	Name:        `E♭ major`,
	Accidentals: -3,
	Type:        Major,
	Notes:       []Note{E_FLAT, F, G, A_FLAT, B_FLAT, C, D},
}

var A_FLAT_MAJOR = Scale{
	id:          AFlatMajor,
	Name:        `A♭ major`,
	Accidentals: -4,
	Type:        Major,
	Notes:       []Note{A_FLAT, B_FLAT, C, D_FLAT, E_FLAT, F, G},
}

// enharmonic to C# major
var D_FLAT_MAJOR = Scale{
	id:          DFlatMajor,
	Name:        `D♭ major`,
	Accidentals: -5,
	Type:        Major,
	Notes:       []Note{D_FLAT, E_FLAT, F, G_FLAT, A_FLAT, B_FLAT, C},
}

// enharmonic to F# major
var G_FLAT_MAJOR = Scale{
	id:          GFlatMajor,
	Name:        `G♭ major`,
	Accidentals: -6,
	Type:        Major,
	Notes:       []Note{G_FLAT, A_FLAT, B_FLAT, C_FLAT, D_FLAT, E_FLAT, F},
}

// enharmonic to E major
var C_FLAT_MAJOR = Scale{
	id:          CFlatMajor,
	Name:        `C♭ major`,
	Accidentals: -7,
	Type:        Major,
	Notes:       []Note{C_FLAT, D_FLAT, E_FLAT, F_FLAT, G_FLAT, A_FLAT, B_FLAT},
}

// Minor scales

var A_MINOR = Scale{
	id:          AMinor,
	Name:        `A minor`,
	Accidentals: 0,
	Type:        Minor,
	Notes:       []Note{A, B, C, D, E, F, G},
}

var E_MINOR = Scale{
	id:          EMinor,
	Name:        `E minor`,
	Accidentals: 1,
	Type:        Minor,
	Notes:       []Note{E, F_SHARP, G, A, B, C, D},
}

var B_MINOR = Scale{
	id:          BMinor,
	Name:        `B minor`,
	Accidentals: 2,
	Type:        Minor,
	Notes:       []Note{B, C_SHARP, D, E, F_SHARP, G, A},
}

var F_SHARP_MINOR = Scale{
	id:          FSharpMinor,
	Name:        `F♯ minor`,
	Accidentals: 3,
	Type:        Minor,
	Notes:       []Note{F_SHARP, G_SHARP, A, B, C_SHARP, D, E},
}

var C_SHARP_MINOR = Scale{
	id:          CSharpMinor,
	Name:        `C♯ minor`,
	Accidentals: 4,
	Type:        Minor,
	Notes:       []Note{C_SHARP, D_SHARP, E, F_SHARP, G_SHARP, A, B},
}

var G_SHARP_MINOR = Scale{
	id:          GSharpMinor,
	Name:        `G♯ minor`,
	Accidentals: 5,
	Type:        Minor,
	Notes:       []Note{G_SHARP, A_SHARP, B, C_SHARP, D_SHARP, E, F_SHARP},
}

// /D♯ minor/E♭ minor
var D_SHARP_MINOR = Scale{
	id:          DSharpMinor,
	Name:        `D♯ minor`,
	Accidentals: 6,
	Type:        Minor,
	Notes:       []Note{D_SHARP, E_SHARP, F_SHARP, G_SHARP, A_SHARP, B, C_SHARP},
}

// A♯ minor/B♭ minor
var A_SHARP_MINOR = Scale{
	id:          ASharpMinor,
	Name:        `A♯ minor`,
	Accidentals: 7,
	Type:        Minor,
	Notes:       []Note{A_SHARP, B_SHARP, C_SHARP, D_SHARP, E_SHARP, F_SHARP, G_SHARP},
}

var D_MINOR = Scale{
	id:          DMinor,
	Name:        `D minor`,
	Accidentals: -1,
	Type:        Minor,
	Notes:       []Note{D, E, F, G, A, B_FLAT, C},
}

var G_MINOR = Scale{
	id:          GMinor,
	Name:        `G minor`,
	Accidentals: -2,
	Type:        Minor,
	Notes:       []Note{G, A, B_FLAT, C, D, E_FLAT, F},
}

var C_MINOR = Scale{
	id:          CMinor,
	Name:        `C minor`,
	Accidentals: -3,
	Type:        Minor,
	Notes:       []Note{C, D, E_FLAT, F, G, A_FLAT, B_FLAT},
}

var F_MINOR = Scale{
	id:          FMinor,
	Name:        `F minor`,
	Accidentals: -4,
	Type:        Minor,
	Notes:       []Note{F, G, A_FLAT, B_FLAT, C, D_FLAT, E_FLAT},
}

// B♭ minor/A♯ minor
var B_FLAT_MINOR = Scale{
	id:          BFlatMinor,
	Name:        `B♭ minor`,
	Accidentals: -5,
	Type:        Minor,
	Notes:       []Note{B_FLAT, C, D_FLAT, E_FLAT, F, G_FLAT, A_FLAT},
}

// E♭ minor/D♯ minor
var E_FLAT_MINOR = Scale{
	id:          EFlatMinor,
	Name:        `E♭ minor`,
	Accidentals: -6,
	Type:        Minor,
	Notes:       []Note{E_FLAT, F, G_FLAT, A_FLAT, B_FLAT, C_FLAT, D_FLAT},
}

// A♭ minor/G♯ minor
var A_FLAT_MINOR = Scale{
	id:          AFlatMinor,
	Name:        `A♭ minor`,
	Accidentals: -7,
	Type:        Minor,
	Notes:       []Note{A_FLAT, B_FLAT, C_FLAT, D_FLAT, E_FLAT, F_FLAT, G_FLAT},
}
