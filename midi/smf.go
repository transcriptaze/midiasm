package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/meta-events"
	"github.com/twystd/midiasm/midi/midi-events"
	"io"
	"sort"
	"time"
)

type SMF struct {
	header *MThd
	tracks []*MTrk
}

type Note struct {
	Channel   byte
	Note      byte
	StartTick uint64
	EndTick   uint64
	Start     time.Duration
	End       time.Duration
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

func (smf *SMF) Render(w io.Writer) {
	smf.header.Render(w)
	for _, track := range smf.tracks {
		track.Render(w)
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

func (smf *SMF) Notes(w io.Writer) error {
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
		for tick, _ := range events {
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

			list := events[tick]
			for _, e := range list {
				if v, ok := e.(*metaevent.Tempo); ok {
					tempo = uint64(v.Tempo)
				}
			}

			for _, e := range list {
				if v, ok := e.(*midievent.NoteOff); ok {
					eventlog.Debug(fmt.Sprintf("NOTE OFF %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t))

					key := uint16(v.Channel)<<8 + uint16(v.Note)
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
				if v, ok := e.(*midievent.NoteOn); ok {
					eventlog.Debug(fmt.Sprintf("NOTE ON  %02X %02X  %-6d %.5f  %s", v.Channel, v.Note, tick, beat, t))

					key := uint16(v.Channel)<<8 + uint16(v.Note)
					note := Note{
						Channel:   v.Channel,
						Note:      v.Note,
						Start:     t,
						StartTick: tick,
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
			fmt.Fprintf(w, "NOTE channel:%d note:%02X start:%-6s end:%-6s\n", n.Channel, n.Note, n.Start, n.End)
		}
	}

	return nil
}
