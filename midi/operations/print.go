package operations

import (
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"os"
	"strings"
	"text/template"
)

const fsmf string = `
>>>>>>>>>>>>>>>>>>>>>>>>>
{{with .MThd}}{{pad (ellipsize .Bytes) 42}}  {{.Tag}} length:{{.Length}}, format:{{.Format}}, tracks:{{.Tracks}}, {{if not .SMPTETimeCode }}metrical time:{{.PPQN}} ppqn{{else}}SMPTE:{{.FPS}} fps,{{.SubFrames}} sub-frames{{end}}{{end}}
{{range .Tracks}}
{{pad (ellipsize .Bytes 0 8)  42}}  {{.Tag}} {{.TrackNumber}} length:{{.Length}}{{range .Events}}
{{pad (ellipsize .Bytes 0 14) 42}}  tick:{{pad .Tick.String 9}}  delta:{{pad .Delta.String 9}}  {{template "events" .}}||{{end}}
{{end}}
>>>>>>>>>>>>>>>>>>>>>>>>>

`
const fevents string = `{{if eq .Tag "TrackName"}}{{template "TrackName" .}}{{else}}   {{pad .Tag 17}}{{end}}`
const ftrackname string = `{{ .Type }} {{pad .Tag 17}}{{ .Name }}`

type Print struct {
	Writer func(midi.Chunk) (io.Writer, error)
}

func (p *Print) Execute(smf *midi.SMF) error {
	functions := template.FuncMap{
		"ellipsize": ellipsize,
		"pad":       pad,
	}

	tmpl, err := template.New("SMF").Funcs(functions).Parse(fsmf)
	if err != nil {
		return err
	}

	_, err = tmpl.New("events").Parse(fevents)
	if err != nil {
		return err
	}

	_, err = tmpl.New("TrackName").Parse(ftrackname)
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, smf)
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

func pad(s string, width int) string {
	if width < 0 {
		return s
	}

	if width < len([]rune(s)) {
		return string([]rune(s)[0:width])
	}

	return s + strings.Repeat(" ", width-len([]rune(s)))
}
