package midi

import (
	"fmt"
	"strings"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

type SMF struct {
	File   string  `json:"-"s`
	MThd   *MThd   `json:"header"`
	Tracks []*MTrk `json:"tracks"`
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

	if smf.MThd.Format == 0 {
		errors = append(errors, smf.validateFormat0()...)
	}

	if smf.MThd.Format == 1 {
		errors = append(errors, smf.validateFormat1()...)
	}

	// End of Track
	for _, track := range smf.Tracks {
		eot := false
		for i, e := range track.Events {
			event := e.Event
			if _, ok := event.(*metaevent.EndOfTrack); ok {
				eot = true
				if i+1 != len(track.Events) {
					errors = append(errors, ValidationError(fmt.Errorf("Track %d: EndOfTrack @%d (%s) is not last event", track.TrackNumber, i+1, clean(event))))
				}
			}
		}

		if !eot {
			errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event", track.TrackNumber)))
		}
	}

	// Program Bank
	for _, track := range smf.Tracks {
		var last interface{}
		for i, e := range track.Events {
			event := e.Event
			c := []*midievent.Controller{nil, nil}

			if cx, ok := last.(*midievent.Controller); ok {
				c[0] = cx
			}

			if cx, ok := event.(*midievent.Controller); ok {
				c[1] = cx
			}

			if c[0] != nil && c[0].Controller.ID == 0x00 && (c[1] == nil || c[1].Controller.ID != 0x20) {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select MSB' event @%d missing LSB (%s)", track.TrackNumber, i, clean(last))))
			}

			if c[1] != nil && c[1].Controller.ID == 0x20 && (c[0] == nil || c[0].Controller.ID != 0x00) {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select LSB' event @%d missing MSB (%s)", track.TrackNumber, i+1, clean(event))))
			}

			if c[0] != nil && c[0].Controller.ID == 0x00 && c[1] != nil && c[1].Controller.ID == 0x20 && c[0].Channel != c[1].Channel {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select MSB' event @%d LSB on another channel (%s)", track.TrackNumber, i, clean(last))))
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select LSB' event @%d MSB on another channel (%s)", track.TrackNumber, i+1, clean(event))))
			}

			last = e.Event
		}
	}

	return errors
}

func (smf *SMF) validateFormat0() []ValidationError {
	errors := []ValidationError{}

	if len(smf.Tracks) != 1 {
		errors = append(errors, ValidationError(fmt.Errorf("File contains %d tracks (expected 1 track for FORMAT 0)", len(smf.Tracks))))
	}

	return errors
}

func (smf *SMF) validateFormat1() []ValidationError {
	errors := []ValidationError{}

	clean := func(e interface{}) string {
		t := fmt.Sprintf("%T", e)
		t = strings.TrimPrefix(t, "*")
		t = strings.TrimPrefix(t, "metaevent.")
		t = strings.TrimPrefix(t, "midievent.")
		t = strings.TrimPrefix(t, "sysex.")

		return t
	}

	if len(smf.Tracks) > 0 {
		track := smf.Tracks[0]
		for _, e := range track.Events {
			event := e.Event
			switch event.(type) {
			case
				*metaevent.Tempo,
				*metaevent.TimeSignature,
				*metaevent.TrackName,
				*metaevent.SMPTEOffset,
				*metaevent.EndOfTrack,
				*metaevent.Copyright:
				continue
			default:
				errors = append(errors, ValidationError(fmt.Errorf("Track 0: unexpected event (%v)", clean(event))))
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

	return errors
}
