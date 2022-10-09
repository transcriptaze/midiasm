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
	"github.com/transcriptaze/midiasm/midi/types"
)

const LOG_TAG = "notes"

type Notes struct {
	Transpose int
	JSON      bool
	Writer    io.Writer
}

type Note struct {
	Channel       types.Channel
	Note          byte
	FormattedNote string
	Velocity      byte
	StartTick     uint64
	EndTick       uint64
	Start         time.Duration
	End           time.Duration
}

func (x *Notes) Execute(smf *midi.SMF) error {
	ppqn := uint64(smf.MThd.Division)
	ctx := context.NewContext()
	tempoMap := make([]*events.Event, 0)

	for _, e := range smf.Tracks[0].Events {
		if _, ok := e.Event.(*metaevent.Tempo); ok {
			tempoMap = append(tempoMap, e)
		}
	}

	for _, track := range smf.Tracks[1:] {
		eventlist := make(map[uint64][]*events.Event, 0)

		for _, e := range tempoMap {
			tick := e.Tick()
			list := eventlist[tick]
			eventlist[tick] = append(list, e)
		}

		for _, e := range track.Events {
			tick := e.Tick()
			list := eventlist[tick]
			eventlist[tick] = append(list, e)
		}

		var ticks []uint64
		for tick, _ := range eventlist {
			ticks = append(ticks, tick)
		}

		sort.SliceStable(ticks, func(i, j int) bool {
			return ticks[i] < ticks[j]
		})

		pending := make(map[uint16]*Note, 0)
		notes := make([]*Note, 0)

		var tempo uint64 = 50000
		var t time.Duration = 0
		var beat float64 = 0.0

		for _, tick := range ticks {
			beat = float64(tick) / float64(ppqn)
			t = time.Duration(1000 * tick * tempo / ppqn)

			if dt := (tick * tempo) % ppqn; dt > 0 {
				warnf("%-5dµs loss of precision converting from tick time to physical time at tick %d", dt, tick)
			}

			list := eventlist[tick]
			for _, e := range list {
				event := e.Event
				if v, ok := event.(*metaevent.Tempo); ok {
					tempo = uint64(v.Tempo)
				}
			}

			for _, e := range list {
				event := e.Event
				if v, ok := event.(*midievent.NoteOff); ok {
					debugf("NOTE OFF %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t)

					key := uint16(v.Channel)<<8 + uint16(v.Note.Value)
					if note := pending[key]; note == nil {
						warnf("NOTE OFF without preceding NOTE ON for %d:%02X", v.Channel, v.Note)
					} else {
						note.End = t
						note.EndTick = tick
						delete(pending, key)
					}
				}
			}

			for _, e := range list {
				event := e.Event
				if v, ok := event.(*metaevent.KeySignature); ok {
					if v.Accidentals < 0 {
						ctx.UseFlats()
					} else {
						ctx.UseSharps()
					}
				}

				if v, ok := event.(*midievent.NoteOn); ok {
					debugf("NOTE ON  %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t)

					key := uint16(v.Channel)<<8 + uint16(v.Note.Value)
					note := Note{
						Channel:       v.Channel,
						Note:          transpose(v.Note.Value, x.Transpose),
						FormattedNote: ctx.FormatNote(transpose(v.Note.Value, x.Transpose)),
						Velocity:      v.Velocity,
						Start:         t,
						StartTick:     tick,
					}

					if pending[key] != nil {
						warnf("NOTE ON without preceding NOTE OFF for %d:%02X", v.Channel, v.Note)
					}

					pending[key] = &note
					notes = append(notes, &note)
				}
			}
		}

		if len(pending) > 0 {
			for k, n := range pending {
				warnf("Incomplete note: %04X %#v", k, n)
			}
		}

		if x.JSON {
			export(notes, x.Writer)
		} else {
			print(notes, x.Writer)
		}
	}

	return nil
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

func print(notes []*Note, w io.Writer) error {
	for _, n := range notes {
		start := n.Start.Truncate(time.Millisecond)
		end := n.End.Truncate(time.Millisecond)
		fmt.Fprintf(w, "%-4s channel:%-2d  note:%02X  velocity:%-3d  start:%-9s  end:%s\n", n.FormattedNote, n.Channel, n.Note, n.Velocity, start, end)
	}

	return nil
}

func export(notes []*Note, w io.Writer) error {
	type note struct {
		Channel  types.Channel `json:"channel"`
		MidiNote byte          `json:"midi-note"`
		Note     string        `json:"note"`
		Velocity byte          `json:"velocity"`
		Start    float64       `json:"start"`
		End      float64       `json:"end"`
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

func debugf(format string, args ...any) {
	log.Debugf(LOG_TAG, format, args...)
}

func warnf(format string, args ...any) {
	log.Warnf(LOG_TAG, format, args...)
}
