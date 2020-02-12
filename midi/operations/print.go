package operations

import (
	"encoding/json"
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

const document string = `{{template "SMF" .}}`

var templates = map[string]string{
	"SMF": `{{template "MThd" .MThd}}
{{range .Tracks}}
{{template "MTrk" .}}{{end}}`,

	"MThd": `{{pad (ellipsize .Bytes 42) 42}}  {{.Tag}} length:{{.Length}}, format:{{.Format}}, tracks:{{.Tracks}}, {{if not .SMPTETimeCode }}metrical time:{{.PPQN}} ppqn{{else}}SMPTE:{{.FPS}} fps,{{.SubFrames}} sub-frames{{end}}`,

	"MTrk": `{{pad (ellipsize .Bytes 24) 42}}  {{.Tag}} {{.TrackNumber}} length:{{.Length}}
{{range .Events}}{{template "event" .}}{{end}}`,

	"event": `{{template "hex" .Bytes}}  tick:{{pad .Tick.String 9}}  delta:{{pad .Delta.String 9}}  {{template "events" .Event}}`,
	"events": `{{if eq .Tag "SequenceNumber"}}{{template "sequenceno" .}}
{{else if eq .Tag "Text"                   }}{{template "text"                   .}}
{{else if eq .Tag "Copyright"              }}{{template "copyright"              .}}
{{else if eq .Tag "TrackName"              }}{{template "trackname"              .}}
{{else if eq .Tag "InstrumentName"         }}{{template "instrumentname"         .}}
{{else if eq .Tag "Lyric"                  }}{{template "lyric"                  .}}
{{else if eq .Tag "Marker"                 }}{{template "marker"                 .}}
{{else if eq .Tag "CuePoint"               }}{{template "cuepoint"               .}}
{{else if eq .Tag "ProgramName"            }}{{template "programname"            .}}
{{else if eq .Tag "DeviceName"             }}{{template "devicename"             .}}
{{else if eq .Tag "MIDIChannelPrefix"      }}{{template "midichannelprefix"      .}}
{{else if eq .Tag "MIDIPort"               }}{{template "midiport"               .}}
{{else if eq .Tag "EndOfTrack"             }}{{template "endoftrack"             .}}
{{else if eq .Tag "Tempo"                  }}{{template "tempo"                  .}}
{{else if eq .Tag "SMPTEOffset"            }}{{template "smpteoffset"            .}}
{{else if eq .Tag "TimeSignature"          }}{{template "timesignature"          .}}
{{else if eq .Tag "KeySignature"           }}{{template "keysignature"           .}}
{{else if eq .Tag "SequencerSpecificEvent" }}{{template "sequencerspecificevent" .}}
{{else if eq .Tag "NoteOff"                }}{{template "noteoff"                .}}
{{else if eq .Tag "NoteOn"                 }}{{template "noteon"                 .}}
{{else if eq .Tag "PolyphonicPressure"     }}{{template "polyphonicpressure"     .}}
{{else if eq .Tag "Controller"             }}{{template "controller"             .}}
{{else if eq .Tag "ProgramChange"          }}{{template "programchange"          .}}
{{else if eq .Tag "ChannelPressure"        }}{{template "channelpressure"        .}}
{{else if eq .Tag "PitchBend"              }}{{template "pitchbend"              .}}
{{else if eq .Tag "SysExMessage"           }}{{template "sysexmessage"           .}}
{{else if eq .Tag "SysExContinuation"      }}{{template "sysexcontinuation"      .}}
{{else if eq .Tag "SysExEscape"            }}{{template "sysexescape"            .}}
{{else                                     }}{{template "unknown"                .}}
{{end}}`,

	"sequenceno":             `{{.Type}} {{pad .Tag 22}} {{.SequenceNumber}}`,
	"text":                   `{{.Type}} {{pad .Tag 22}} {{.Text}}`,
	"copyright":              `{{.Type}} {{pad .Tag 22}} {{.Copyright}}`,
	"trackname":              `{{.Type}} {{pad .Tag 22}} {{.Name}}`,
	"instrumentname":         `{{.Type}} {{pad .Tag 22}} {{.Name}}`,
	"lyric":                  `{{.Type}} {{pad .Tag 22}} {{.Lyric}}`,
	"marker":                 `{{.Type}} {{pad .Tag 22}} {{.Marker}}`,
	"cuepoint":               `{{.Type}} {{pad .Tag 22}} {{.CuePoint}}`,
	"programname":            `{{.Type}} {{pad .Tag 22}} {{.Name}}`,
	"devicename":             `{{.Type}} {{pad .Tag 22}} {{.Name}}`,
	"midichannelprefix":      `{{.Type}} {{pad .Tag 22}} {{.Channel}}`,
	"midiport":               `{{.Type}} {{pad .Tag 22}} {{.Port}}`,
	"endoftrack":             `{{.Type}} {{    .Tag   }}`,
	"tempo":                  `{{.Type}} {{pad .Tag 22}} tempo:{{.Tempo}}`,
	"smpteoffset":            `{{.Type}} {{pad .Tag 22}} {{.Hour}} {{.Minute}} {{.Second}} {{.FrameRate}} {{.Frames}} {{.FractionalFrames}}`,
	"timesignature":          `{{.Type}} {{pad .Tag 22}} {{.Numerator}}/{{.Denominator}}, {{.TicksPerClick }} ticks per click, {{.ThirtySecondsPerQuarter}}/32 per quarter`,
	"keysignature":           `{{.Type}} {{pad .Tag 22}} {{.Key}}`,
	"sequencerspecificevent": `{{.Type}} {{pad .Tag 22}} {{.Manufacturer.Name}}, {{.Data}}`,

	"noteoff":            `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"noteon":             `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"polyphonicpressure": `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} pressure:{{.Pressure}}`,
	"controller":         `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} controller:{{.Controller}}, value:{{.Value}}`,
	"programchange":      `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} program:{{.Program }}`,
	"channelpressure":    `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} pressure:{{.Pressure}}`,
	"pitchbend":          `{{.Status}} {{pad .Tag 22}} channel:{{pad .Channel 2}} bend:{{.Bend}}`,

	"sysexmessage":      `{{.Status}} {{pad .Tag 22}} {{.Manufacturer.Name}}, {{.Data}}`,
	"sysexcontinuation": `{{.Status}} {{pad .Tag 22}} {{.Data}}`,
	"sysexescape":       `{{.Status}} {{pad .Tag 22}} {{.Data}}`,

	"unknown": `?? {{.Tag}}`,
	"hex":     `{{pad (ellipsize (valign . 3) 42) 42}}`,
}

type Print struct {
	root *template.Template
}

func NewPrint() (*Print, error) {
	functions := template.FuncMap{
		"ellipsize": ellipsize,
		"pad":       pad,
		"valign":    valign,
	}

	tmpl, err := template.New("document").Funcs(functions).Parse(document)
	if err != nil {
		return nil, err
	}

	for name, t := range templates {
		if _, err = tmpl.New(name).Parse(t); err != nil {
			return nil, err
		}
	}

	return &Print{
		root: tmpl,
	}, nil
}

func (p *Print) LoadTemplates(r io.Reader) error {
	templates := struct {
		Templates map[string]string `json:"templates"`
	}{
		Templates: make(map[string]string, 0),
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &templates)
	if err != nil {
		return err
	}

	for name, t := range templates.Templates {
		if _, err = p.root.New(name).Parse(t); err != nil {
			return err
		}
	}

	return nil
}

func (p *Print) Print(object interface{}, template string, w io.Writer) error {
	tmpl := p.root.Lookup(template)
	if tmpl == nil {
		return fmt.Errorf("'%s' does not match any defined template", template)
	}

	return tmpl.Execute(w, object)
}

func ellipsize(v interface{}, length int) string {
	if length <= 0 {
		return ""
	}

	s := []rune(fmt.Sprintf("%v", v))
	if len(s) <= length {
		return string(s)
	}

	return string(s[0:length-1]) + `…`
}

func pad(v interface{}, width int) string {
	s := fmt.Sprintf("%v", v)
	if width < len([]rune(s)) {
		return s
	}

	return s + strings.Repeat(" ", width-len([]rune(s)))
}

func valign(bytes types.Hex, w ...int) string {
	ix := 0

	for {
		b := bytes[ix]
		ix += 1

		if b&0x80 == 0 {
			break
		}
	}

	width := 4
	if len(w) > 0 {
		width = w[0]
	}

	pad := ""
	if width > ix {
		pad = strings.Repeat("   ", width-ix)
	}

	return fmt.Sprintf("%s%v", pad, bytes)
}
