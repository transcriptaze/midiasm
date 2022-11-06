package assemble

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/encoding/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

type TextAssembler struct {
}

func NewTextAssembler() TextAssembler {
	return TextAssembler{}
}

func (a TextAssembler) Assemble(r io.Reader) ([]byte, error) {
	chunks, err := a.read(r)
	if err != nil {
		return nil, err
	}

	// ... parse chunks

	smf := midi.SMF{}

	for _, chunk := range chunks {
		for _, line := range chunk {
			switch {
			case strings.HasPrefix(line, "%%"):
				// comment - ignore

			case strings.Contains(line, "MThd"):
				if mthd, err := a.parseMThd(chunk); err != nil {
					return nil, err
				} else {
					smf.MThd = mthd
				}

				break

			case strings.Contains(line, "MTrk"):
				if mtrk, err := a.parseMTrk(chunk); err != nil {
					return nil, err
				} else {
					smf.MThd.Tracks += 1
					smf.Tracks = append(smf.Tracks, mtrk)
				}

				break
			}
		}
	}

	// ... 'k, done

	var b bytes.Buffer
	var e = midifile.NewEncoder(&b)

	if err := e.Encode(smf); err != nil {
		return nil, err
	} else {
		return b.Bytes(), nil
	}
}

func (a TextAssembler) read(r io.Reader) ([][]string, error) {
	scanner := bufio.NewScanner(r)
	lines := make(chan string)
	chunks := make(chan []string)

	go a.scan(scanner, lines)
	go a.chunkify(lines, chunks)

	list := [][]string{}
	for chunk := range chunks {
		list = append(list, chunk)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (a TextAssembler) scan(scanner *bufio.Scanner, lines chan string) {
	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)
}

func (a TextAssembler) chunkify(lines chan string, chunks chan []string) {
	tags := regexp.MustCompile("(MThd)|(MTrk)")

	clone := func(slice []string) []string {
		chunk := make([]string, len(slice))
		copy(chunk, slice)
		return chunk
	}

	var chunk []string
	for line := range lines {
		if strings.Contains(line, "MThd") {
			chunk = []string{line}
			break
		}
	}

	for line := range lines {
		if match := tags.FindStringSubmatch(line); match != nil {
			chunks <- clone(chunk)
			chunk = []string{line}
		} else {
			chunk = append(chunk, line)
		}
	}

	if len(chunk) > 0 {
		chunks <- clone(chunk)
	}

	close(chunks)
}

func (a TextAssembler) parseMThd(chunk []string) (*midi.MThd, error) {
	for _, line := range chunk {
		if strings.Contains(line, "MThd") {
			var format uint16
			var ppqn uint16

			if match := regexp.MustCompile(`format:(0|1|2)`).FindStringSubmatch(line); match == nil || len(match) < 2 {
				return nil, fmt.Errorf("missing or invalid 'format' field in MThd")
			} else if v, err := strconv.ParseUint(match[1], 10, 16); err != nil {
				return nil, err
			} else {
				format = uint16(v)
			}

			if match := regexp.MustCompile(`metrical(?:[ -])?time:([0-9]+)\s*ppqn`).FindStringSubmatch(line); match == nil || len(match) < 2 {
				return nil, fmt.Errorf("missing 'metrical-time' field in MThd")
			} else if v, err := strconv.ParseUint(match[1], 10, 16); err != nil {
				return nil, err
			} else {
				ppqn = uint16(v)
			}

			if mthd, err := midi.NewMThd(format, 0, ppqn); err != nil {
				return nil, err
			} else {
				return mthd, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid MThd")
}

func (a TextAssembler) parseMTrk(chunk []string) (*midi.MTrk, error) {
	lines := make(chan string)
	closed := make(chan bool, 1)

	defer func() {
		closed <- true
	}()

	go func() {
		for _, line := range chunk {
			select {
			case lines <- line:
			case <-closed:
				break
			}
		}

		close(lines)
	}()

	// ... make MTrk
	var mtrk *midi.MTrk

	for line := range lines {
		if strings.Contains(line, "MTrk") {
			if v, err := midi.NewMTrk(); err != nil {
				return nil, err
			} else {
				mtrk = v
				break
			}
		}
	}

	if mtrk == nil {
		return nil, fmt.Errorf("missing MTrk")
	}

	// ... extract events
	type E encoding.TextUnmarshaler

	f := func(line string, e E) error {
		if err := e.UnmarshalText([]byte(line)); err != nil {
			return err
		} else {
			mtrk.Events = append(mtrk.Events, events.NewEvent(0, 0, e, nil))
		}

		return nil
	}

	g := map[string]func() E{
		"SequenceNumber":         func() E { return &metaevent.SequenceNumber{} },
		"Text":                   func() E { return &metaevent.Text{} },
		"Copyright":              func() E { return &metaevent.Copyright{} },
		"TrackName":              func() E { return &metaevent.TrackName{} },
		"InstrumentName":         func() E { return &metaevent.InstrumentName{} },
		"Lyric":                  func() E { return &metaevent.Lyric{} },
		"Marker":                 func() E { return &metaevent.Marker{} },
		"CuePoint":               func() E { return &metaevent.CuePoint{} },
		"ProgramName":            func() E { return &metaevent.ProgramName{} },
		"DeviceName":             func() E { return &metaevent.DeviceName{} },
		"MIDIChannelPrefix":      func() E { return &metaevent.MIDIChannelPrefix{} },
		"MIDIPort":               func() E { return &metaevent.MIDIPort{} },
		"Tempo":                  func() E { return &metaevent.Tempo{} },
		"TimeSignature":          func() E { return &metaevent.TimeSignature{} },
		"KeySignature":           func() E { return &metaevent.KeySignature{} },
		"SMPTEOffset":            func() E { return &metaevent.SMPTEOffset{} },
		"EndOfTrack":             func() E { return &metaevent.EndOfTrack{} },
		"SequencerSpecificEvent": func() E { return &metaevent.SequencerSpecificEvent{} },
		"ProgramChange":          func() E { return &midievent.ProgramChange{} },
		"Controller":             func() E { return &midievent.Controller{} },
		"NoteOn":                 func() E { return &midievent.NoteOn{} },
		"NoteOff":                func() E { return &midievent.NoteOff{} },
		"PolyphonicPressure":     func() E { return &midievent.PolyphonicPressure{} },
		"ChannelPressure":        func() E { return &midievent.ChannelPressure{} },
		"PitchBend":              func() E { return &midievent.PitchBend{} },
	}

	for line := range lines {
		if strings.HasPrefix(line, "%%") {
			continue
		}

		for k, v := range g {
			if strings.Contains(line, k) {
				e := v()
				if err := f(line, e); err != nil {
					return nil, err
				} else if _, ok := e.(*metaevent.EndOfTrack); ok {
					return mtrk, nil
				}
			}
		}
	}

	return mtrk, fmt.Errorf("missing EndOfTrack")
}
