package disassemble

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/transcriptaze/midiasm/midi/types"
)

//go:embed template
var document string

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
