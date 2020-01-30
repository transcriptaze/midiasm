package context

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
)

var Sharps = map[byte]string{
	0:  "C",
	1:  "C\u266f",
	2:  "D",
	3:  "D\u266f",
	4:  "E",
	5:  "F",
	6:  "F\u266f",
	7:  "G",
	8:  "G\u266f",
	9:  "A",
	10: "A\u266f",
	11: "B",
}

var Flats = map[byte]string{
	0:  "C",
	1:  "D\u266d",
	2:  "D",
	3:  "E\u266d",
	4:  "E",
	5:  "F",
	6:  "G\u266d",
	7:  "G",
	8:  "A\u266d",
	9:  "A",
	10: "B\u266d",
	11: "B",
}

type Context struct {
	scale map[byte]string
	casio bool
	notes map[uint16]string
}

func NewContext() *Context {
	return &Context{
		scale: Sharps,
		casio: false,
		notes: make(map[uint16]string),
	}
}

func (ctx *Context) Scale() map[byte]string {
	return ctx.scale
}

func (ctx *Context) UseSharps() *Context {
	ctx.scale = Sharps

	return ctx
}

func (ctx *Context) UseFlats() *Context {
	ctx.scale = Flats

	return ctx
}

func (ctx *Context) Casio() bool {
	return ctx.casio
}

func (ctx *Context) CasioOn() *Context {
	ctx.casio = true

	return ctx
}

func (ctx *Context) CasioOff() *Context {
	ctx.casio = false

	return ctx
}

func (ctx *Context) FormatNote(n byte) string {
	note := ctx.scale[n%12]
	octave := int(n/12) - 2

	return fmt.Sprintf("%s%d", note, octave)
}

func (ctx *Context) GetNoteOff(ch types.Channel, n byte) string {
	key := uint16(ch)
	key <<= 8
	key |= uint16(n)

	if note, ok := ctx.notes[key]; ok {
		return note
	}

	return ctx.FormatNote(n)
}

func (ctx *Context) PutNoteOn(ch types.Channel, n byte) {
	key := uint16(ch)
	key <<= 8
	key |= uint16(n)

	ctx.notes[key] = ctx.FormatNote(n)
}
