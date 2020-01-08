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
{{pad (ellipsize .MThd.Bytes 0 14) 42}}  {{template "MThd" .MThd}}
{{range .Tracks}}
{{pad (ellipsize      .Bytes 0 8)  42}}  {{template "MTrk" .}}{{range .Events}}
{{pad (ellipsize      .Bytes 0 14) 42}}  {{template "event" .}}{{end}}
{{end}}
>>>>>>>>>>>>>>>>>>>>>>>>>

`
const fMThd string = `{{.Tag}} length:{{.Length}}, format:{{.Format}}, tracks:{{.Tracks}}, {{if not .SMPTETimeCode }}metrical time:{{.PPQN}} ppqn{{else}}SMPTE:{{.FPS}} fps,{{.SubFrames}} sub-frames{{end}}`
const fMTrk string = `{{.Tag}} {{.TrackNumber}} length:{{.Length}}`
const fEvent string = `tick:{{pad .Tick.String 9}}  delta:{{pad .Delta.String 9}}  {{template "events" .}}||`
const fEvents string = `{{if eq .Tag "TrackName"}}{{template "TrackName" .}}{{else}}XX {{pad .Tag 17}}{{end}}`
const fTrackName string = `{{ .Type }} {{pad .Tag 17}}{{ .Name }}`

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

	_, err = tmpl.New("MThd").Parse(fMThd)
	if err != nil {
		return err
	}

	_, err = tmpl.New("MTrk").Parse(fMTrk)
	if err != nil {
		return err
	}

	_, err = tmpl.New("event").Parse(fEvent)
	if err != nil {
		return err
	}

	_, err = tmpl.New("events").Parse(fEvents)
	if err != nil {
		return err
	}

	_, err = tmpl.New("TrackName").Parse(fTrackName)
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
