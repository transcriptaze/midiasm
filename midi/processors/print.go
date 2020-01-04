package processors

import (
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"os"
	"strings"
	"text/template"
)

type Print struct {
	Writer func(midi.Chunk) (io.Writer, error)
}

func (p *Print) Execute(smf *midi.SMF) error {
	format := `
>>>>>>>>>>>>>>>>>>>>>>>>>
{{ellipsize .MThd.Bytes 42}}  {{.MThd.Tag}} length:{{.MThd.Length}}, format:{{.MThd.Format}}, tracks:{{.MThd.Tracks}}, metrical time:{{.MThd.PPQN}} ppqn"
{{range .Tracks}}
{{ellipsize .Bytes 42 0 8}}  {{.Tag}} {{.TrackNumber}} length:{{.Length}}
{{range .Events}}
{{ellipsize .Bytes 42}}  tick:{{.Tick}}   delta:{{.Delta}}{{end}}
{{end}}
>>>>>>>>>>>>>>>>>>>>>>>>>

`

	functions := template.FuncMap{
		"ellipsize": func(bytes types.Hex, width int, offsets ...int) string {
			start := 0
			end := len(bytes)

			if len(offsets) > 0 && offsets[0] > 0 {
				start = offsets[0]
			}

			if len(offsets) > 1 && offsets[1] > start {
				end = offsets[1]
			}

			hex := bytes[start:end].String()
			if end-start < len(bytes) {
				hex += `…`
			}

			if width < 0 {
				return hex
			}

			if width == 0 {
				return ""
			}

			if width < len([]rune(hex)) {
				return string([]rune(hex)[0:width])
			}

			return hex + strings.Repeat(" ", width-len([]rune(hex)))
		},
	}

	tmpl, err := template.New("SMF").Funcs(functions).Parse(format)
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
