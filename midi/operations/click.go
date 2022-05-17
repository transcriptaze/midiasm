package operations

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/twystd/midiasm/midi"
	// "github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	// "github.com/twystd/midiasm/midi/events/midi"
	"github.com/twystd/midiasm/midi/types"
)

type ClickTrack struct {
	Writer io.Writer
}

type Click struct {
	Channel       types.Channel
	Note          byte
	FormattedNote string
	Velocity      byte
	StartTick     uint64
	EndTick       uint64
	Start         time.Duration
	End           time.Duration
}

func (x *ClickTrack) Execute(smf *midi.SMF) error {
	ppqn := uint64(smf.MThd.Division)
	// ctx := context.NewContext()
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
			if list == nil {
				list = make([]*events.Event, 0)
			}

			eventlist[tick] = append(list, e)
		}

		for _, e := range track.Events {
			if _, ok := e.Event.(*metaevent.TimeSignature); ok {
				tick := uint64(e.Tick)
				list := eventlist[tick]
				if list == nil {
					list = make([]*events.Event, 0)
				}

				eventlist[tick] = append(list, e)
			}
		}

		// for _, v := range eventlist {
		// 	for _, e := range v {
		// 		tick := uint64(e.Tick)
		// 		fmt.Printf(">>> time signature: %v %v\n", tick, e.Event)
		// 	}
		// }

		var ticks []uint64
		for tick, _ := range eventlist {
			ticks = append(ticks, tick)
		}

		sort.SliceStable(ticks, func(i, j int) bool {
			return ticks[i] < ticks[j]
		})

		changes := make([]*metaevent.TimeSignature, 0)

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
				if v, ok := event.(*metaevent.TimeSignature); ok {
					eventlog.Debug(fmt.Sprintf("TIME SIGNATURE %d/%d  %-6d %.5f  %s", v.Numerator, v.Denominator, tick, beat, t))
					changes = append(changes, v)
				}
			}
		}

		for _, t := range changes {
			fmt.Fprintf(x.Writer, "%d/%d\n", t.Numerator, t.Denominator)
		}
	}

	return fmt.Errorf("NOT IMPLEMENTED")
}
