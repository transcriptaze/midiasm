package click

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
)

type ClickTrack struct {
	Writer io.Writer
}

type Cluck struct {
	Bar           uint
	Tempo         uint
	TimeSignature string
}

type signature struct {
	numerator   uint8
	denominator uint8
}

func (x *ClickTrack) Execute(smf *midi.SMF) error {
	clucks := map[uint]Cluck{}
	ppqn := uint64(smf.MThd.Division)
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
			if _, ok := e.Event.(*metaevent.TimeSignature); ok {
				tick := uint64(e.Tick)
				list := eventlist[tick]
				eventlist[tick] = append(list, e)
			}
		}

		var ticks []uint64
		for tick, _ := range eventlist {
			ticks = append(ticks, tick)
		}

		sort.SliceStable(ticks, func(i, j int) bool {
			return ticks[i] < ticks[j]
		})

		var tempo uint = 120
		var beat float64 = 0.0
		var bar uint = 1
		var last float64 = 0.0
		var timeSignature *metaevent.TimeSignature

		for _, tick := range ticks {
			beat = float64(tick) / float64(ppqn)

			list := eventlist[tick]
			for _, e := range list {
				event := e.Event
				if v, ok := event.(*metaevent.Tempo); ok {
					tempo = uint(math.Round(60.0 * 1000000.0 / float64(v.Tempo)))

					eventlog.Debug(fmt.Sprintf("TEMPO %v  bar:%-3d", v, bar))

					cluck := clucks[bar]
					cluck.Bar = bar
					cluck.Tempo = tempo
					clucks[bar] = cluck
				}
			}

			for _, e := range list {
				event := e.Event
				if v, ok := event.(*metaevent.TimeSignature); ok {
					if timeSignature == nil {
						last = beat
					} else {
						bar += uint((beat - last) / float64(timeSignature.Numerator))
						last = beat
					}

					eventlog.Debug(fmt.Sprintf("TIME SIGNATURE %v  bar:%-3d", v, bar))

					cluck := clucks[bar]
					cluck.Bar = bar
					cluck.Tempo = tempo
					cluck.TimeSignature = fmt.Sprintf("%v", v)
					clucks[bar] = cluck

					timeSignature = v
				}
			}
		}

		// Only process Tracks 0 and 1 for now
		break
	}

	summary := []Cluck{}
	for _, v := range clucks {
		summary = append(summary, v)
	}

	sort.SliceStable(summary, func(i, j int) bool {
		return summary[i].Bar < summary[j].Bar
	})

	for _, v := range summary {
		fmt.Fprintf(x.Writer, "bar %-4v  tempo:%-3v  time-signature %v\n", v.Bar, v.Tempo, v.TimeSignature)
	}

	return fmt.Errorf("NOT IMPLEMENTED")
}
