package operations

import (
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"os"
	"strings"
	"text/template"
)

const document string = `
>>>>>>>>>>>>>>>>>>>>>>>>>
{{pad (ellipsize .MThd.Bytes 0 14) 42}}  {{template "MThd" .MThd}}

{{range .Tracks}}{{pad (ellipsize      .Bytes 0 8)  42}}  {{template "MTrk" .}}
{{range .Events}}{{pad (ellipsize      .Bytes 0 14) 42}}  {{template "event" .}}{{end}}
{{end}}
>>>>>>>>>>>>>>>>>>>>>>>>>

`

var templates = map[string]string{
	"MThd": `{{.Tag}} length:{{.Length}}, format:{{.Format}}, tracks:{{.Tracks}}, {{if not .SMPTETimeCode }}metrical time:{{.PPQN}} ppqn{{else}}SMPTE:{{.FPS}} fps,{{.SubFrames}} sub-frames{{end}}`,
	"MTrk": `{{.Tag}} {{.TrackNumber}} length:{{.Length}}`,

	"event": `tick:{{pad .Tick.String 9}}  delta:{{pad .Delta.String 9}}  {{template "events" .}}`,
	"events": `{{if eq .Tag "TrackName"}}{{template "trackname" .}}
{{else if eq .Tag "EndOfTrack"         }}{{template "endoftrack"         .}}
{{else if eq .Tag "Tempo"              }}{{template "tempo"              .}}
{{else if eq .Tag "TimeSignature"      }}{{template "timesignature"      .}}
{{else if eq .Tag "KeySignature"       }}{{template "keysignature"       .}}
{{else if eq .Tag "NoteOff"            }}{{template "noteoff"            .}}
{{else if eq .Tag "NoteOn"             }}{{template "noteon"             .}}
{{else if eq .Tag "PolyphonicPressure" }}{{template "polyphonicpressure" .}}
{{else if eq .Tag "Controller"         }}{{template "controller"         .}}
{{else if eq .Tag "ProgramChange"      }}{{template "programchange"      .}}
{{else if eq .Tag "ChannelPressure"    }}{{template "channelpressure"    .}}
{{else                                 }}XX {{pad .Tag 16}} 
{{end}}`,

	"trackname":     `{{.Type}} {{pad .Tag 18}} {{.Name}}`,
	"endoftrack":    `{{.Type}} {{    .Tag   }}`,
	"tempo":         `{{.Type}} {{pad .Tag 18}} tempo:{{.Tempo }}`,
	"timesignature": `{{.Type}} {{pad .Tag 18}} {{.Numerator}}/{{.Denominator}}, {{.TicksPerClick }} ticks per click, {{.ThirtySecondsPerQuarter}}/32 per quarter`,
	"keysignature":  `{{.Type}} {{pad .Tag 18}} {{.Key }}`,

	"noteoff":            `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"noteon":             `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"polyphonicpressure": `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} pressure:{{.Pressure}}`,
	"controller":         `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} controller:{{.Controller}}, value:{{.Value}}`,
	"programchange":      `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} program:{{.Program }}`,
	"channelpressure":    `{{.Status}} {{pad .Tag 18}} channel:{{pad .Channel 2}} pressure:{{.Pressure}}`,
}

type Print struct {
	Writer func(midi.Chunk) (io.Writer, error)
}

func (p *Print) Execute(smf *midi.SMF) error {
	err := p.printWithTemplate(smf, os.Stdout)
	if err != nil {
		return err
	}

	if w, err := p.Writer(smf.MThd); err != nil {
		return err
	} else {
		smf.MThd.Print(w)
	}

	for _, track := range smf.Tracks {
		if w, err := p.Writer(track); err != nil {
			return err
		} else {
			track.Print(w)
		}
	}

	return nil
}

func (p *Print) printWithTemplate(smf *midi.SMF, w io.Writer) error {
	functions := template.FuncMap{
		"ellipsize": ellipsize,
		"pad":       pad,
	}

	tmpl, err := template.New("SMF").Funcs(functions).Parse(document)
	if err != nil {
		return err
	}

	for name, t := range templates {
		if _, err = tmpl.New(name).Parse(t); err != nil {
			return err
		}
	}

	return tmpl.Execute(w, smf)
}

func ellipsize(bytes types.Hex, offsets ...int) string {
	start := 0
	end := len(bytes)

	if len(offsets) > 0 && offsets[0] > 0 {
		start = offsets[0]
	}

	if len(offsets) > 1 && offsets[1] > start && offsets[1] < end {
		end = offsets[1]
	}

	hex := bytes[start:end].String()
	if end-start < len(bytes) {
		hex += `…`
	}

	return hex
}

func pad(v interface{}, width int) string {
	s := fmt.Sprintf("%v", v)
	if width < len([]rune(s)) {
		return s
	}

	return s + strings.Repeat(" ", width-len([]rune(s)))
}
