package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type SMF struct {
	File   string
	MThd   *MThd
	Tracks []*MTrk
}

type Note struct {
	Channel       byte
	Note          byte
	FormattedNote string
	Velocity      byte
	StartTick     uint64
	EndTick       uint64
	Start         time.Duration
	End           time.Duration
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

	mthd, ok := chunks[0].(*MThd)
	if !ok {
		return fmt.Errorf("invalid MIDI file - expected MThd chunk, got %T", chunks[0])
	}

	if len(chunks[1:]) != int(mthd.Tracks) {
		return fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", mthd.Tracks, len(chunks[1:]))
	}

	tracks := make([]*MTrk, 0)
	for i, chunk := range chunks[1:] {
		track, ok := chunk.(*MTrk)
		if !ok {
			return fmt.Errorf("invalid MIDI file - expected MTrk chunk, got %T", chunk)
		}

		track.TrackNumber = types.TrackNumber(i)
		tracks = append(tracks, track)
	}

	smf.MThd = mthd
	smf.Tracks = tracks

	return nil
}

func (smf *SMF) LoadConfiguration(r io.Reader) error {
	conf := struct {
		Manufacturers []types.Manufacturer `json:"manufacturers"`
	}{
		Manufacturers: make([]types.Manufacturer, 0),
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return err
	}

	for _, m := range conf.Manufacturers {
		if err := types.AddManufacturer(m); err != nil {
			return err
		}
	}

	return nil
}

func (smf *SMF) Validate() []ValidationError {
	errors := []ValidationError{}

	clean := func(e events.IEvent) string {
		t := fmt.Sprintf("%T", e)
		t = strings.TrimPrefix(t, "*")
		t = strings.TrimPrefix(t, "metaevent.")
		t = strings.TrimPrefix(t, "midievent.")
		t = strings.TrimPrefix(t, "sysex.")

		return t
	}

	if smf.MThd.Format == 0 && len(smf.Tracks) != 1 {
		errors = append(errors, ValidationError(fmt.Errorf("File contains %d tracks (expected 1 track for FORMAT 0)", len(smf.Tracks))))
	}

	if smf.MThd.Format == 1 {
		if len(smf.Tracks) > 0 {
			track := smf.Tracks[0]
			for _, event := range track.Events {
				switch event.(type) {
				case *metaevent.Tempo,
					*metaevent.TrackName,
					*metaevent.EndOfTrack:
					continue
				default:
					errors = append(errors, ValidationError(fmt.Errorf("Track 0: unexpected event (%s)", clean(event))))
				}
			}
		}

		for _, track := range smf.Tracks[1:] {
			for _, event := range track.Events {
				switch event.(type) {
				case *metaevent.Tempo:
					errors = append(errors, ValidationError(fmt.Errorf("Track %d: unexpected event (%s)", track.TrackNumber, clean(event))))
				}
			}
		}
	}

	for _, track := range smf.Tracks {
		if len(track.Events) == 0 {
			errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event", track.TrackNumber)))
		} else {
			event := track.Events[len(track.Events)-1]
			if _, ok := event.(*metaevent.EndOfTrack); !ok {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event (%s)", track.TrackNumber, clean(event))))
			}
		}
	}

	return errors
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
		err := mtrk.UnmarshalBinary(bytes)
		if err != nil {
			return nil, err
		}
		return &mtrk, nil
	}

	return nil, nil
}
