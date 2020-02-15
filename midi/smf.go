package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/types"
	"io"
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

	// Extract header
	chunk, err := readChunk(r)
	if err == nil && chunk != nil {
		var mthd MThd
		if err := mthd.UnmarshalBinary(chunk); err == nil {
			smf.MThd = &mthd
		}
	}

	// Extract tracks
	for err == nil {
		chunk, err = readChunk(r)
		if err == nil && chunk != nil {
			if string(chunk[0:4]) == "MTrk" {
				mtrk := MTrk{
					TrackNumber: types.TrackNumber(len(smf.Tracks)),
				}

				if err := mtrk.UnmarshalBinary(chunk); err != nil {
					return err
				}

				smf.Tracks = append(smf.Tracks, &mtrk)
			}
		}
	}

	if err != io.EOF {
		return err
	}

	if len(smf.Tracks) != int(smf.MThd.Tracks) {
		return fmt.Errorf("number of tracks in file does not match MThd - expected %d, got %d", smf.MThd.Tracks, len(smf.Tracks))
	}

	return nil
}

func (smf *SMF) Validate() []ValidationError {
	errors := []ValidationError{}

	clean := func(e interface{}) string {
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
			for _, e := range track.Events {
				event := e.Event
				switch event.(type) {
				case *metaevent.Tempo,
					*metaevent.TrackName,
					*metaevent.SMPTEOffset,
					*metaevent.EndOfTrack:
					continue
				default:
					errors = append(errors, ValidationError(fmt.Errorf("Track 0: unexpected event (%s)", clean(event))))
				}
			}
		}

		for _, track := range smf.Tracks[1:] {
			for _, e := range track.Events {
				event := e.Event
				switch event.(type) {
				case *metaevent.Tempo:
					errors = append(errors, ValidationError(fmt.Errorf("Track %d: unexpected event (%s)", track.TrackNumber, clean(event))))

				case *metaevent.SMPTEOffset:
					errors = append(errors, ValidationError(fmt.Errorf("Track %d: unexpected event (%s)", track.TrackNumber, clean(event))))
				}
			}
		}
	}

	for _, track := range smf.Tracks {
		if len(track.Events) == 0 {
			errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event", track.TrackNumber)))
		} else {
			e := track.Events[len(track.Events)-1]
			event := e.Event
			if _, ok := event.(*metaevent.EndOfTrack); !ok {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event (%s)", track.TrackNumber, clean(event))))
			}
		}
	}

	return errors
}

func readChunk(r *bufio.Reader) ([]byte, error) {
	peek, err := r.Peek(8)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(peek[4:8])
	bytes := make([]byte, length+8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
