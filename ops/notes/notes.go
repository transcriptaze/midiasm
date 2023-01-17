package notes

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/transcriptaze/midiasm/log"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/lib"
)

const LOG_TAG = "notes"

type Notes struct {
	Transpose int
	JSON      bool
	Writer    io.Writer
}

type Note struct {
	Channel       lib.Channel
	Note          byte
	FormattedNote string
	Velocity      byte
	StartTick     uint64
	EndTick       uint64
	Start         time.Duration
	End           time.Duration
}

type event struct {
	tick  uint64
	delta uint64
	at    time.Duration
	event any
	ctx   context.Context
}

func (x *Notes) Execute(smf *midi.SMF) error {
	if notes, err := extract(smf, x.Transpose); err != nil {
		return err
	} else {
		if x.JSON {
			export(notes, x.Writer)
		} else {
			print(notes, x.Writer)
		}
	}

	return nil
}

func extract(smf *midi.SMF, transposition int) ([]Note, error) {
	notes := make([]Note, 0)

	// ... build tempo map
	tempoMap := buildTempoMap(*smf)

	// ... extract track notes
	for _, track := range smf.Tracks[1:] {
		if events, err := buildTrackEvents(*track, tempoMap, smf.MThd.PPQN); err != nil {
			return nil, err
		} else if list, err := buildNoteList(events, transposition); err != nil {
			return nil, err
		} else {
			notes = append(notes, list...)
		}
	}

	return notes, nil
}

func buildTempoMap(smf midi.SMF) []events.Event {
	list := []events.Event{}

	for _, e := range smf.Tracks[0].Events {
		if _, ok := e.Event.(metaevent.Tempo); ok {
			list = append(list, *e)
		}
	}

	return list
}

//			if dt := (tick * tempo) % ppqn; dt > 0 {
//				warnf("%-5dµs loss of precision converting from tick time to physical time at tick %d", dt, tick)
//			}

func buildTrackEvents(track midi.MTrk, tempoMap []events.Event, ppqn uint16) ([]event, error) {
	ctx := context.NewContext()

	// ... build event list
	list := []event{}
	for _, e := range track.Events {
		list = append(list, event{
			tick:  e.Tick(),
			delta: uint64(e.Delta()),
			event: e.Event,
		})
	}

	for _, e := range tempoMap {
		list = append(list, event{
			tick:  e.Tick(),
			delta: 0,
			event: e.Event,
		})
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].tick < list[j].tick
	})

	// ... assign event timestamps
	tempo := uint64(50000)
	at := time.Duration(0)

	for i, e := range list {
		delta := time.Duration(1000 * e.delta * tempo / uint64(ppqn))
		at += delta

		if v, ok := e.event.(metaevent.Tempo); ok {
			tempo = uint64(v.Tempo)
		}

		if v, ok := e.event.(metaevent.KeySignature); ok {
			if v.Accidentals < 0 {
				ctx.UseFlats()
			} else {
				ctx.UseSharps()
			}
		}

		list[i].at = at.Round(1 * time.Millisecond)
		list[i].ctx = *ctx
	}

	// ... extract NoteOn and NoteOff events
	notes := []event{}
	for i := range list {
		e := list[i]
		switch e.event.(type) {
		case midievent.NoteOn:
			notes = append(notes, e)

		case midievent.NoteOff:
			notes = append(notes, e)

		case metaevent.EndOfTrack:
			break
		}
	}

	return notes, nil
}

func buildNoteList(events []event, transposition int) ([]Note, error) {
	notes := []Note{}

	pending := map[lib.Channel]map[uint8]event{}
	for i := 0; i < 16; i++ {
		pending[lib.Channel(i)] = map[uint8]event{}
	}

	for i := range events {
		e := events[i]
		if v, ok := e.event.(midievent.NoteOn); ok && v.Velocity > 0 {
			if _, ok := pending[v.Channel][v.Note.Value]; !ok {
				pending[v.Channel][v.Note.Value] = e
			} else {
				warnf("NoteOn without preceding NoteOff for %d:%02X", v.Channel, v.Note)
			}
		}

		if v, ok := e.event.(midievent.NoteOn); ok && v.Velocity == 0 {
			if p, ok := pending[v.Channel][v.Note.Value]; ok {
				on := p.event.(midievent.NoteOn)
				notes = append(notes, Note{
					Channel:       on.Channel,
					Note:          transpose(on.Note.Value, transposition),
					FormattedNote: midievent.FormatNote(&p.ctx, transpose(on.Note.Value, transposition)),
					Velocity:      on.Velocity,
					Start:         p.at,
					StartTick:     p.tick,
					EndTick:       e.tick,
					End:           e.at,
				})

				delete(pending[v.Channel], v.Note.Value)
			} else {
				warnf("NoteOff without preceding NoteOn for %d:%02X", v.Channel, v.Note)
			}
		}

		if v, ok := e.event.(midievent.NoteOff); ok {
			debugf(e)
			if p, ok := pending[v.Channel][v.Note.Value]; ok {
				on := p.event.(midievent.NoteOn)
				notes = append(notes, Note{
					Channel:       on.Channel,
					Note:          transpose(on.Note.Value, transposition),
					FormattedNote: midievent.FormatNote(&p.ctx, transpose(on.Note.Value, transposition)),
					Velocity:      on.Velocity,
					Start:         p.at,
					StartTick:     p.tick,
					EndTick:       e.tick,
					End:           e.at,
				})

				delete(pending[v.Channel], v.Note.Value)
			} else {
				warnf("NoteOff without preceding NoteOn for %d:%02X", v.Channel, v.Note)
			}
		}
	}

	for i := 0; i < 16; i++ {
		list := pending[lib.Channel(i)]
		if len(list) > 0 {
			for k, n := range list {
				warnf("Incomplete note: %04X %#v", k, n)
			}
		}
	}

	return notes, nil
}

func transpose(note uint8, transpose int) uint8 {
	v := int(note) + transpose

	if v < 0 {
		return 0
	} else if v > 255 {
		return 255
	} else {
		return uint8(v)
	}
}

func print(notes []Note, w io.Writer) error {
	for _, n := range notes {
		start := n.Start.Truncate(time.Millisecond)
		end := n.End.Truncate(time.Millisecond)
		duration := (n.End - n.Start).Truncate(time.Millisecond)

		fmt.Fprintf(w, "%-4s channel:%-2d  note:%02X  velocity:%-3d  start:%-9s  end:%-9s  duration:%s\n", n.FormattedNote, n.Channel, n.Note, n.Velocity, start, end, duration)
	}

	return nil
}

func export(notes []Note, w io.Writer) error {
	type note struct {
		Channel  lib.Channel `json:"channel"`
		MidiNote byte        `json:"midi-note"`
		Note     string      `json:"note"`
		Velocity byte        `json:"velocity"`
		Start    float64     `json:"start"`
		End      float64     `json:"end"`
	}

	object := struct {
		Notes []note `json:"notes"`
	}{}

	for _, n := range notes {
		start := n.Start.Truncate(time.Millisecond)
		end := n.End.Truncate(time.Millisecond)

		object.Notes = append(object.Notes, note{
			Channel:  n.Channel,
			Note:     n.FormattedNote,
			MidiNote: n.Note,
			Velocity: n.Velocity,
			Start:    start.Seconds(),
			End:      end.Seconds(),
		})
	}

	if bytes, err := json.MarshalIndent(object, "", "  "); err != nil {
		return err
	} else {
		w.Write(bytes)
	}

	return nil
}

func debugf(e event) {
	fmt := "%-8s %02X %02X %-3v %-3v %-6d %s"

	switch v := e.event.(type) {
	case midievent.NoteOn:
		log.Debugf(LOG_TAG, fmt, "NOTE ON", v.Channel, v.Note.Value, v.Note.Name, v.Velocity, e.tick, e.at)

	case midievent.NoteOff:
		log.Debugf(LOG_TAG, fmt, "NOTE OFF", v.Channel, v.Note.Value, v.Note.Name, v.Velocity, e.tick, e.at)
	}
}

func warnf(format string, args ...any) {
	log.Warnf(LOG_TAG, format, args...)
}
