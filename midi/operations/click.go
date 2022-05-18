package operations

import (
	"fmt"
	"io"
	"sort"

	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
)

type ClickTrack struct {
	Writer io.Writer
}

// type Click struct {
// 	Channel       types.Channel
// 	Note          byte
// 	FormattedNote string
// 	Velocity      byte
// 	StartTick     uint64
// 	EndTick       uint64
// 	Start         time.Duration
// 	End           time.Duration
// }

type signature struct {
	numerator   uint8 
	denominator uint8
}

func (x *ClickTrack) Execute(smf *midi.SMF) error {
	type change struct {
		Bar       uint32
		Signature signature
	}

	changes := []change{}
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

		var tempo uint64 = 50000
		// var t time.Duration = 0
		var beat float64 = 0.0
		var bar uint32 = 0
		var last float64 = 0.0
		var timeSignature *metaevent.TimeSignature

		for _, tick := range ticks {
			beat = float64(tick) / float64(ppqn)
			// t = time.Duration(1000 * tick * tempo / ppqn)

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
				if v, ok := event.(*metaevent.TimeSignature); ok {
					if timeSignature == nil {
						bar = 1
						last = beat
					} else {
						bar += uint32((beat - last) / float64(timeSignature.Numerator))
						last = beat
					}

					eventlog.Debug(fmt.Sprintf("TIME SIGNATURE %d/%d  bar:%-3d", v.Numerator, v.Denominator, bar))
					changes = append(changes, change{
						Bar: bar,
						Signature: signature{
							numerator:   v.Numerator,
							denominator: v.Denominator,
						},
					})

					timeSignature = v
				}
			}
		}

		// Only process Tracks 0 and 1 for now
		break
	}

	for _, v := range changes {
		fmt.Fprintf(x.Writer, "bar %-4v  time-signature %d/%d\n", v.Bar, v.Signature.numerator, v.Signature.denominator)
	}

	return fmt.Errorf("NOT IMPLEMENTED")
}
