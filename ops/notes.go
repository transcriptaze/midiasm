package ops

import (
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/events/midi"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"sort"
	"time"
)

type Notes struct {
	Writer io.Writer
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
			tick := uint64(e.Tick)
			list := eventlist[tick]
			eventlist[tick] = append(list, e)
		}

		for _, e := range track.Events {
			tick := uint64(e.Tick)
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
				eventlog.Warn(fmt.Sprintf("%-5dµs loss of precision converting from tick time to physical time at tick %d", dt, tick))
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
					eventlog.Debug(fmt.Sprintf("NOTE OFF %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t))

					key := uint16(v.Channel)<<8 + uint16(v.Note.Value)
					if note := pending[key]; note == nil {
						eventlog.Warn(fmt.Sprintf("NOTE OFF without preceding NOTE ON for %d:%02X", v.Channel, v.Note))
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
					eventlog.Debug(fmt.Sprintf("NOTE ON  %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t))

					key := uint16(v.Channel)<<8 + uint16(v.Note.Value)
					note := Note{
						Channel:       v.Channel,
						Note:          v.Note.Value,
						FormattedNote: ctx.FormatNote(v.Note.Value),
						Velocity:      v.Velocity,
						Start:         t,
						StartTick:     tick,
					}

					if pending[key] != nil {
						eventlog.Warn(fmt.Sprintf("NOTE ON without preceding NOTE OFF for %d:%02X", v.Channel, v.Note))
					}

					pending[key] = &note
					notes = append(notes, &note)
				}
			}
		}

		if len(pending) > 0 {
			for k, n := range pending {
				eventlog.Warn(fmt.Sprintf("Incomplete note: %04X %#v", k, n))
			}
		}

		for _, n := range notes {
			start := n.Start.Truncate(time.Millisecond)
			end := n.End.Truncate(time.Millisecond)
			fmt.Fprintf(x.Writer, "%-4s channel:%-2d  note:%02X  velocity:%-3d  start:%-9s  end:%s\n", n.FormattedNote, n.Channel, n.Note, n.Velocity, start, end)
		}
	}

	return nil
}
