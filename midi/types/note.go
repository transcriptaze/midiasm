package types

type Note struct {
	ord  int
	Name string
}

func (n Note) String() string {
	return n.Name
}

var C = Note{
	ord:  0,
	Name: `C`,
}

var C_SHARP = Note{
	ord:  1,
	Name: `C♯`,
}

var D_FLAT = Note{
	ord:  1,
	Name: `D♭`,
}

var D = Note{
	ord:  2,
	Name: `D`,
}

var D_SHARP = Note{
	ord:  3,
	Name: `D♯`,
}

var E_FLAT = Note{
	ord:  3,
	Name: `E♭`,
}

var E = Note{
	ord:  4,
	Name: `E`,
}

var E_SHARP = Note{
	ord:  5,
	Name: `E♯`,
}

var F_FLAT = Note{
	ord:  4,
	Name: `F♭`,
}

var F = Note{
	ord:  5,
	Name: `F`,
}

var F_SHARP = Note{
	ord:  6,
	Name: `F♯`,
}

var G_FLAT = Note{
	ord:  6,
	Name: `G♭`,
}

var G = Note{
	ord:  7,
	Name: `G`,
}

var G_SHARP = Note{
	ord:  8,
	Name: `G♯`,
}

var A_FLAT = Note{
	ord:  8,
	Name: `A♭`,
}

var A = Note{
	ord:  9,
	Name: `A`,
}

var A_SHARP = Note{
	ord:  10,
	Name: `A♯`,
}

var B_FLAT = Note{
	ord:  10,
	Name: `B♭`,
}

var B = Note{
	ord:  11,
	Name: `B`,
}

var B_SHARP = Note{
	ord:  0,
	Name: `B♯`,
}

var C_FLAT = Note{
	ord:  11,
	Name: `C♭`,
}
