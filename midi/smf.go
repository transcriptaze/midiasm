package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"github.com/twystd/midiasm/midi/meta-events"
	"github.com/twystd/midiasm/midi/midi-events"
	"io"
	"os"
	"sort"
	"time"
)

type SMF struct {
	header *MThd
	tracks []*MTrk
}

func (smf *SMF) UnmarshalBinary(data []byte) error {
	r := bufio.NewReader(bytes.NewReader(data))
	chunks := make([]Chunk, 0)

	chunk := Chunk(nil)
	err := error(nil)

	for err == nil {
		chunk, err = readChunk(r)
		if err == nil && chunk != nil {
			chunks = append(chunks, chunk)
		}
	}

	if err != io.EOF {
		return err
	}

	if len(chunks) == 0 {
		return fmt.Errorf("contains no MIDI chunks")
	}

	header, ok := chunks[0].(*MThd)
	if !ok {
		return fmt.Errorf("invalid MIDI file - expected MThd chunk, got %T", chunks[0])
	}

	if len(chunks[1:]) != int(header.tracks) {
		return fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", header.tracks, len(chunks[1:]))
	}

	tracks := make([]*MTrk, len(chunks[1:]))
	for i, chunk := range chunks[1:] {
		if track, ok := chunk.(*MTrk); ok {
			tracks[i] = track
		} else {
			return fmt.Errorf("invalid MIDI file - expected MTrk chunk, got %T", chunk)
		}
	}

	smf.header = header
	smf.tracks = tracks

	return nil
}

func (smf *SMF) Render() {
	smf.header.Render(os.Stdout)
	for _, track := range smf.tracks {
		track.Render(os.Stdout)
	}
}

func readChunk(r *bufio.Reader) (Chunk, error) {
	peek, err := r.Peek(8)
	if err != nil {
		return nil, err
	}

	tag := string(peek[0:4])
	length := binary.BigEndian.Uint32(peek[4:8])

	bytes := make([]byte, length+8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}

	switch tag {
	case "MThd":
		var mthd MThd
		if err := mthd.UnmarshalBinary(bytes); err != nil {
			return nil, err
		}
		return &mthd, nil

	case "MTrk":
		var mtrk MTrk
		if err := mtrk.UnmarshalBinary(bytes); err != nil {
			return nil, err
		}
		return &mtrk, nil
	}

	return nil, nil
}

func (smf *SMF) Notes() error {
	w := os.Stdout
	ppqn := uint64(smf.header.division)
	tempoMap := make([]event.IEvent, 0)

	for _, e := range smf.tracks[0].Events {
		if v, ok := e.(*metaevent.Tempo); ok {
			tempoMap = append(tempoMap, v)
		}
	}

	for _, track := range smf.tracks[1:] {
		events := make(map[uint64][]event.IEvent, 0)

		for _, e := range tempoMap {
			tick := e.TickValue()
			list := events[tick]
			if list == nil {
				list = make([]event.IEvent, 0)
			}

			events[tick] = append(list, e)
		}

		for _, e := range track.Events {
			tick := e.TickValue()
			list := events[tick]
			if list == nil {
				list = make([]event.IEvent, 0)
			}

			events[tick] = append(list, e)
		}

		var ticks []uint64
		for tickv, _ := range events {
			ticks = append(ticks, tickv)
		}

		sort.SliceStable(ticks, func(i, j int) bool {
			return ticks[i] < ticks[j]
		})

		var tempo uint64 = 50000
		var tick uint64 = 0
		var t time.Duration = 0
		var beat float64 = 0.0

		for _, tickv := range ticks {
			tick = tickv
			beat = float64(tick) / float64(ppqn)
			t = time.Duration(1000 * tick * tempo / ppqn)

			if dt := (tick * tempo) % ppqn; dt > 0 {
				fmt.Printf("WARNING: %dµs loss of precision converting from tick time to physical time at tick %d\n", dt, tick)
			}

			list := events[tickv]
			for _, e := range list {
				if v, ok := e.(*metaevent.Tempo); ok {
					tempo = uint64(v.Tempo)
				}
			}

			for _, e := range list {
				if v, ok := e.(*midievent.NoteOff); ok {
					fmt.Fprintf(w, "NOTE OFF %02X  %-6d %.5f  %-10d %.5f\n", v.Note, tick, beat, t.Microseconds(), t.Seconds())
				}
			}

			for _, e := range list {
				if v, ok := e.(*midievent.NoteOn); ok {
					fmt.Fprintf(w, "NOTE ON  %02X  %-6d %.5f  %-10d %.5f\n", v.Note, tick, beat, t.Microseconds(), t.Seconds())
				}
			}
		}

		fmt.Fprintln(w)
	}

	return nil
}
