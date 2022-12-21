package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/commands"
	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type tsv struct {
	out     string
	conf    string
	c4      bool
	verbose bool
	debug   bool
}

var TSV = tsv{
	c4:      false,
	verbose: false,
	debug:   false,
}

func (t tsv) GetCommand() (string, commands.Command) {
	return "tsv", TSV
}

func (t tsv) Flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("tsv", flag.ExitOnError)

	flagset.StringVar(&TSV.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&TSV.c4, "C4", TSV.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&TSV.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&TSV.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (t tsv) Help() {
	fmt.Println()
	fmt.Println("  Extracts the MIDI information as TSV for use with e.g. a spreadsheet.")
	fmt.Println()
	fmt.Println("    midiasm tsv [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      <MIDI file>  MIDI file to export as JSON.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the TSV to a file. Default is to write to stdout.")
	fmt.Println("      --C4          Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug       Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose     Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm tsv --debug --verbose --out one-time.tsv one-time.mid")
	fmt.Println()
}

func (t tsv) MiddleC() lib.MiddleC {
	if t.c4 {
		return lib.C4
	}

	return lib.C3
}

func (t tsv) Config() string {
	return t.conf
}

func (t tsv) Debug() bool {
	return t.debug
}

func (t tsv) Verbose() bool {
	return t.verbose
}

func (t tsv) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

	if smf, err := decode(filename); err != nil {
		return err
	} else if err := validate(smf); err != nil {
		return err
	} else {
		return export(smf)
	}
}

func decode(filename string) (*midi.SMF, error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer f.Close()

		decoder := midifile.NewDecoder()

		if smf, err := decoder.Decode(f); err != nil {
			return nil, err
		} else if smf == nil {
			return nil, fmt.Errorf("failed to decode MIDI file")
		} else {
			return smf, nil
		}
	}
}

func validate(smf *midi.SMF) error {
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

func export(smf *midi.SMF) error {
	// ... TSV header record
	header := []string{}
	for range smf.Tracks {
		header = append(header, []string{"Tick", "Delta", "Tag", "Channel", "Details"}...)
	}

	// ... build track columns
	tracks := [][][]string{}

	for _, t := range smf.Tracks {
		track := [][]string{}
		for _, e := range t.Events {
			track = append(track, fields(e.Event))
		}

		tracks = append(tracks, track)
	}

	// ... zip tracks
	rows := 0
	for _, track := range tracks {
		if len(track) > rows {
			rows = len(track)
		}
	}

	records := make([][]string, rows)

	for i, _ := range records {
		record := []string{}
		for _, track := range tracks {
			if i < len(track) {
				record = append(record, track[i]...)
			} else {
				record = append(record, []string{"", "", ""}...)
			}
		}

		records[i] = record
	}

	// ... export as TSV
	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

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
		details = fmt.Sprintf("%v|%v, %v", e.Note.Value, e.Note.Name, e.Velocity)
	}

	if e, ok := v.(midievent.NoteOn); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v|%v, %v", e.Note.Value, e.Note.Name, e.Velocity)
	}

	if e, ok := v.(midievent.PolyphonicPressure); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v", e.Pressure)
	}

	if e, ok := v.(midievent.Controller); ok {
		channel = fmt.Sprintf("%v", e.Channel)
		details = fmt.Sprintf("%v|%v, %v", e.Controller.ID, e.Controller.Name, e.Value)
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
		details = fmt.Sprintf("%v|%v|%v, %v", e.Manufacturer.ID, e.Manufacturer.Region, e.Manufacturer.Name, e.Data)
	}

	if e, ok := v.(sysex.SysExContinuationMessage); ok {
		details = fmt.Sprintf("%v", e.Data)
	}

	if e, ok := v.(sysex.SysExEscapeMessage); ok {
		details = fmt.Sprintf("%v", e.Data)
	}
	return []string{tick, delta, tag, channel, details}
}
