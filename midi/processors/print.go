package processors

import (
	"github.com/twystd/midiasm/midi"
	"io"
	"os"
	"text/template"
)

type Print struct {
	Writer func(midi.Chunk) (io.Writer, error)
}

func (p *Print) Execute(smf *midi.SMF) error {
	format := `
>>>>>>>>>>>>>>>>>>>>>>>>>
{{.MThd.Bytes}}   {{.MThd.Tag}} length:{{.MThd.Length}}, format:{{.MThd.Format}}, tracks:{{.MThd.Tracks}}, metrical time:{{.MThd.PPQN}} ppqn"
{{range .Tracks}}
{{slice .Bytes 0 8}}…                    {{.Tag}} {{.TrackNumber}} length:{{.Length}}{{end}}
>>>>>>>>>>>>>>>>>>>>>>>>>

`
	tmpl, err := template.New("SMF").Parse(format)
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
