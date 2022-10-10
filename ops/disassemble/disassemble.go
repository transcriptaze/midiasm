package disassemble

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/transcriptaze/midiasm/midi/types"
)

const document string = `{{template "SMF" .}}`

var templates = map[string]string{
	"SMF": `{{template "MThd" .MThd}}
{{range .Tracks}}
{{template "MTrk" .}}{{end}}`,

	"MThd": `{{pad 42 (ellipsize .Bytes 42) }}  {{.Tag}} length:{{.Length}}, format:{{.Format}}, tracks:{{.Tracks}}, {{if not .SMPTETimeCode }}metrical time:{{.PPQN}} ppqn{{else}}SMPTE:{{.FPS}} fps,{{.SubFrames}} sub-frames{{end}}`,

	"MTrk": `{{pad 42 (ellipsize .Bytes 24) }}  {{.Tag}} {{.TrackNumber}} length:{{.Length}}
{{range .Events}}{{template "event" .}}{{end}}`,

	"event": `{{template "hex" .Bytes}}  tick:{{.Tick | pad 9}}  delta:{{pad 9 .Delta}}  {{template "events" .Event}}`,
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

	"sequenceno":             `{{.Type}} {{pad 22 .Tag}} {{.SequenceNumber}}`,
	"text":                   `{{.Type}} {{pad 22 .Tag}} {{.Text}}`,
	"copyright":              `{{.Type}} {{pad 22 .Tag}} {{.Copyright}}`,
	"trackname":              `{{.Type}} {{pad 22 .Tag}} {{.Name}}`,
	"instrumentname":         `{{.Type}} {{pad 22 .Tag}} {{.Name}}`,
	"lyric":                  `{{.Type}} {{pad 22 .Tag}} {{.Lyric}}`,
	"marker":                 `{{.Type}} {{pad 22 .Tag}} {{.Marker}}`,
	"cuepoint":               `{{.Type}} {{pad 22 .Tag}} {{.CuePoint}}`,
	"programname":            `{{.Type}} {{pad 22 .Tag}} {{.Name}}`,
	"devicename":             `{{.Type}} {{pad 22 .Tag}} {{.Name}}`,
	"midichannelprefix":      `{{.Type}} {{pad 22 .Tag}} {{.Channel}}`,
	"midiport":               `{{.Type}} {{pad 22 .Tag}} {{.Port}}`,
	"endoftrack":             `{{.Type}} {{    .Tag   }}`,
	"tempo":                  `{{.Type}} {{pad 22 .Tag}} tempo:{{.Tempo}}`,
	"smpteoffset":            `{{.Type}} {{pad 22 .Tag}} {{.Hour}} {{.Minute}} {{.Second}} {{.FrameRate}} {{.Frames}} {{.FractionalFrames}}`,
	"timesignature":          `{{.Type}} {{pad 22 .Tag}} {{.Numerator}}/{{.Denominator}}, {{.TicksPerClick }} ticks per click, {{.ThirtySecondsPerQuarter}}/32 per quarter`,
	"keysignature":           `{{.Type}} {{pad 22 .Tag}} {{.Key}}`,
	"sequencerspecificevent": `{{.Type}} {{pad 22 .Tag}} {{.Manufacturer.Name}}, {{.Data}}`,

	"noteoff":            `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"noteon":             `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} note:{{.Note.Name}}, velocity:{{.Velocity}}`,
	"polyphonicpressure": `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} pressure:{{.Pressure}}`,
	"controller":         `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} {{.Controller.ID}}/{{.Controller.Name}}, value:{{.Value}}`,
	"programchange":      `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} bank:{{.Bank}}, program:{{.Program }}`,
	"channelpressure":    `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} pressure:{{.Pressure}}`,
	"pitchbend":          `{{.Status}} {{pad 22 .Tag}} channel:{{pad 2 .Channel}} bend:{{.Bend}}`,

	"sysexmessage":      `{{.Status}} {{pad 22 .Tag}} {{.Manufacturer.Name}}, {{.Data}}`,
	"sysexcontinuation": `{{.Status}} {{pad 22 .Tag}} {{.Data}}`,
	"sysexescape":       `{{.Status}} {{pad 22 .Tag}} {{.Data}}`,

	"unknown": `?? {{.Tag}}`,
	"hex":     `{{pad 42 (ellipsize (valign . 3) 42) }}`,
}

type Disassemble struct {
	root *template.Template
}

func NewDisassemble() (*Disassemble, error) {
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

	return &Disassemble{
		root: tmpl,
	}, nil
}

func (p *Disassemble) LoadTemplates(r io.Reader) error {
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

func (p *Disassemble) Print(smf any, template string, w io.Writer) error {
	tmpl := p.root.Lookup(template)
	if tmpl == nil {
		return fmt.Errorf("'%s' does not match any defined template", template)
	}

	return tmpl.Execute(w, smf)
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

func pad(width int, v any) string {
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
