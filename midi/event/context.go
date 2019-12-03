package event

import (
	"fmt"
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
	Scale map[byte]string
}

func (ctx *Context) FormatNote(n byte) string {
	note := ctx.Scale[n%12]
	octave := int(n/12) - 2

	return fmt.Sprintf("%s%d", note, octave)
}
