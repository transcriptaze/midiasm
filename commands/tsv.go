package commands

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
)

type tsv struct {
	out       string
	delimiter string
	tabular   bool
}

var TSV = tsv{}

func (t *tsv) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&TSV.out, "out", "", "Output file path (or directory for split files)")
	flagset.StringVar(&TSV.delimiter, "delimiter", "", "Column delimiter.Defaults to TAB")
	flagset.BoolVar(&TSV.tabular, "tabular", false, "Formats the output as fixed width columns")

	return flagset
}

func (t tsv) Help() {
	fmt.Println()
	fmt.Println("  Extracts the MIDI information as TSV for use with e.g. a spreadsheet.")
	fmt.Println()
	fmt.Println("    midiasm tsv [--debug] [--verbose] [--C4] [--out <file>] [--delimiter <string>] <MIDI file>")
	fmt.Println()
	fmt.Println("      <MIDI file>  MIDI file to export as JSON.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --out <file>          Writes the TSV to a file. Default is to write to stdout.")
	fmt.Println("      --delimiter <string>  Column separator (defaults to TAB).")
	fmt.Println("      --tabular             Formats the output as fixed width columns.")
	fmt.Println("      --C4                  Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug               Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose             Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm tsv --debug --verbose --out one-time.tsv one-time.mid")
	fmt.Println()
}

func (t tsv) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

	if b, err := os.ReadFile(filename); err != nil {
		return err
	} else if smf, err := t.decode(bytes.NewBuffer(b)); err != nil {
		return err
	} else if err := t.validate(smf); err != nil {
		return err
	} else if header, records, err := t.export(smf); err != nil {
		return err
	} else {
		if t.out == "" {
			return writeTable(records, os.Stdout)
		} else if f, err := os.Create(t.out); err != nil {
			return err
		} else {
			defer f.Close()

			if t.tabular {
				return writeTable(records, f)
			} else if t.delimiter == "" || t.delimiter == `\t` {
				return writeTSV(header, records, '\t', f)
			} else {
				return writeTSV(header, records, []rune(t.delimiter)[0], f)
			}
		}
	}
}

func (t tsv) export(smf *midi.SMF) ([]string, [][]string, error) {
	// ... build table
	header := []string{"Tick", "Delta", "Tag", "Channel", "Note", "Velocity", "Details"}
	columns := 0
	rows := 0
	for _, t := range smf.Tracks {
		columns += len(header)
		if len(t.Events) > rows {
			rows = len(t.Events)
		}
	}

	records := make([][]string, rows+1)
	for i := range records {
		records[i] = make([]string, columns)
	}

	// ... header record
	for i := range smf.Tracks {
		offset := i * len(header)
		for j, h := range header {
			col := offset + j
			records[0][col] = h
		}
	}

	// ... build track columns
	for i, t := range smf.Tracks {
		offset := i * len(header)
		for j, e := range t.Events {
			row := j + 1
			list := fields(e.Event)

			for k, f := range list {
				col := offset + k
				records[row][col] = f
			}
		}
	}

	return header, records, nil
}

func (t *tsv) decode(r io.Reader) (*midi.SMF, error) {
	decoder := midifile.NewDecoder()

	if smf, err := decoder.Decode(r); err != nil {
		return nil, err
	} else if smf == nil {
		return nil, fmt.Errorf("failed to decode MIDI file")
	} else {
		return smf, nil
	}
	// }
}

func (t *tsv) validate(smf *midi.SMF) error {
	errors := smf.Validate()

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "WARNING: there are validation errors:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "         ** %v\n", e)
		}
		fmt.Fprintln(os.Stderr)
	}

	return nil
}

func writeTable(records [][]string, w io.Writer) error {
	columns := 0
	for _, record := range records {
		if len(record) > columns {
			columns = len(record)
		}
	}

	widths := make([]int, columns)

	for _, record := range records {
		for i, f := range record {
			if len(f) > widths[i] {
				widths[i] = len(f)
			}
		}
	}

	formats := make([]string, len(widths))
	for i, w := range widths {
		formats[i] = fmt.Sprintf("%%-%vv", w)
	}

	for _, record := range records {
		row := []string{}
		for i, f := range record {
			row = append(row, fmt.Sprintf(formats[i], f))
		}

		fmt.Fprintf(w, "%v\n", strings.Join(row, "  "))
	}

	return nil
}

func writeTSV(header []string, records [][]string, delimiter rune, out io.Writer) error {
	w := csv.NewWriter(out)
	w.Comma = delimiter

	if err := w.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func fields(v events.IEvent) []string {
	tick := fmt.Sprintf("%v", v.Tick())
	delta := fmt.Sprintf("%v", v.Delta())
	tag := v.Tag()
	channel := ""
	note := ""
	velocity := ""
	details := ""

	// ... meta events
	if e, ok := v.(metaevent.SequenceNumber); ok {
		details = fmt.Sprintf("%v", e.SequenceNumber)
	}

	if e, ok := v.(metaevent.Text); ok {
		details = fmt.Sprintf("%v", e.Text)
	}

	if e, ok := v.(metaevent.Copyright); ok {
		details = fmt.Sprintf("%v", e.Copyright)
	}

	if e, ok := v.(metaevent.TrackName); ok {
		details = fmt.Sprintf("%v", e.Name)
	}

	if e, ok := v.(metaevent.InstrumentName); ok {
		details = fmt.Sprintf("%v", e.Name)
	}

	if e, ok := v.(metaevent.Lyric); ok {
		details = fmt.Sprintf("%v", e.Lyric)
	}

	if e, ok := v.(metaevent.Marker); ok {
		details = fmt.Sprintf("%v", e.Marker)
	}

	if e, ok := v.(metaevent.CuePoint); ok {
		details = fmt.Sprintf("%v", e.CuePoint)
	}

	if e, ok := v.(metaevent.ProgramName); ok {
		details = fmt.Sprintf("%v", e.Name)
	}

	if e, ok := v.(metaevent.DeviceName); ok {
		details = fmt.Sprintf("%v", e.Name)
	}

	if e, ok := v.(metaevent.MIDIChannelPrefix); ok {
		details = fmt.Sprintf("%v", e.Channel)
	}

	if e, ok := v.(metaevent.MIDIPort); ok {
		details = fmt.Sprintf("%v", e.Port)
	}

	if _, ok := v.(metaevent.EndOfTrack); ok {
	}

	if e, ok := v.(metaevent.Tempo); ok {
		details = fmt.Sprintf("%v", e.Tempo)
	}

	if e, ok := v.(metaevent.SMPTEOffset); ok {
		details = fmt.Sprintf("offset %02v:%02v:%02v  frame rate %v:%v:%v", e.Hour, e.Minute, e.Second, e.FrameRate, e.Frames, e.FractionalFrames)
	}

	if e, ok := v.(metaevent.KeySignature); ok {
		details = fmt.Sprintf("%v", e.Key)
	}

	if e, ok := v.(metaevent.TimeSignature); ok {
		details = fmt.Sprintf("%v/%v, %v ticks/click, %v 32nds/quarter", e.Numerator, e.Denominator, e.TicksPerClick, e.ThirtySecondsPerQuarter)
	}

	if e, ok := v.(metaevent.SequencerSpecificEvent); ok {
		details = fmt.Sprintf("%v [%v]", e.Manufacturer.Name, e.Data)
	}

	// ... MIDI events
	if e, ok := v.(midievent.NoteOff); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		note = fmt.Sprintf("%v", e.Note.Value)
		velocity = fmt.Sprintf("%v", e.Velocity)
		details = fmt.Sprintf("%v", e.Note.Name)
	}

	if e, ok := v.(midievent.NoteOn); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		note = fmt.Sprintf("%v", e.Note.Value)
		velocity = fmt.Sprintf("%v", e.Velocity)
		details = fmt.Sprintf("%v", e.Note.Name)
	}

	if e, ok := v.(midievent.PolyphonicPressure); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v", e.Pressure)
	}

	if e, ok := v.(midievent.Controller); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v:%v, %v", e.Controller.ID, e.Controller.Name, e.Value)
	}

	if e, ok := v.(midievent.ProgramChange); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v, %v", e.Bank, e.Program)
	}

	if e, ok := v.(midievent.ChannelPressure); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v", e.Pressure)
	}

	if e, ok := v.(midievent.PitchBend); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v", e.Bend)
	}

	// ... SysEx events
	if e, ok := v.(sysex.SysExMessage); ok {
		details = fmt.Sprintf("%v:%v:%v, %v", e.Manufacturer.ID, e.Manufacturer.Region, e.Manufacturer.Name, e.Data)
	}

	if e, ok := v.(sysex.SysExContinuationMessage); ok {
		details = fmt.Sprintf("%v", e.Data)
	}

	if e, ok := v.(sysex.SysExEscapeMessage); ok {
		details = fmt.Sprintf("%v", e.Data)
	}

	return []string{tick, delta, tag, channel, note, velocity, details}
}
