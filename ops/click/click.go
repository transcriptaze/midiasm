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

type Span struct {
	Start         uint
	End           uint
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

			if _, ok := e.Event.(*metaevent.EndOfTrack); ok {
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
			if timeSignature != nil {
				bar += uint((beat - last) / float64(timeSignature.Numerator))
			}
			last = beat

			list := eventlist[tick]
			for _, e := range list {
				event := e.Event

				// ... tempo changes
				if v, ok := event.(*metaevent.Tempo); ok {
					tempo = uint(math.Round(60.0 * 1000000.0 / float64(v.Tempo)))

					eventlog.Debug(fmt.Sprintf("%-14v  %-5v bar:%-3d", "TEMPO", v, bar))

					cluck := clucks[bar]
					cluck.Bar = bar
					cluck.Tempo = tempo
					clucks[bar] = cluck
				}

				// ... time signature changes
				if v, ok := event.(*metaevent.TimeSignature); ok {
					eventlog.Debug(fmt.Sprintf("%-14v  %-5v bar:%-3d", "TIME SIGNATURE", v, bar))

					cluck := clucks[bar]
					cluck.Bar = bar
					cluck.Tempo = tempo
					cluck.TimeSignature = fmt.Sprintf("%v", v)
					clucks[bar] = cluck

					timeSignature = v
				}

				if _, ok := event.(*metaevent.EndOfTrack); ok {
					eventlog.Debug(fmt.Sprintf("%-14v  %-5v bar:%-3d", "END OF TRACK", "", bar))

					cluck := clucks[bar]
					cluck.Bar = bar - 1 // because end of track as after the last bar ????
					cluck.Tempo = tempo
					cluck.TimeSignature = fmt.Sprintf("%v", timeSignature)
					clucks[bar] = cluck
				}
			}
		}

		// Only process Track 1 for now
		break
	}

	list := []Cluck{}
	for _, v := range clucks {
		list = append(list, v)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Bar < list[j].Bar
	})

	for _, v := range list {
		fmt.Fprintf(x.Writer, "bar %-4v  tempo:%-3v  time-signature %v\n", v.Bar, v.Tempo, v.TimeSignature)
	}

	spans := []Span{}
	if len(list) > 0 {
		span := Span{
			Start:         list[0].Bar,
			End:           list[0].Bar,
			Tempo:         list[0].Tempo,
			TimeSignature: list[0].TimeSignature,
		}

		for _, v := range list[1:] {
			span.End = v.Bar - 1
			spans = append(spans, span)

			span = Span{
				Start:         v.Bar,
				End:           v.Bar,
				Tempo:         v.Tempo,
				TimeSignature: v.TimeSignature,
			}
		}
	}

	fmt.Println()
	for _, v := range spans {
		fmt.Fprintf(x.Writer, "bars %-7v  tempo:%-3v  time-signature %v\n", fmt.Sprintf("%v:%v", v.Start, v.End), v.Tempo, v.TimeSignature)
	}
	fmt.Println()

	return fmt.Errorf("NOT IMPLEMENTED")
}
