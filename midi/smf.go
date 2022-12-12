package midi

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

type SMF struct {
	MThd   *MThd   `json:"header"`
	Tracks []*MTrk `json:"tracks"`
}

func (smf *SMF) Validate() []ValidationError {
	errors := []ValidationError{}

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
			if events.IsEndOfTrack(e) {
				eot = true
				if i+1 != len(track.Events) {
					errors = append(errors, ValidationError(fmt.Errorf("Track %d: EndOfTrack @%d (%s) is not last event", track.TrackNumber, i+1, events.Clean(e))))
				}
			}
		}

		if !eot {
			errors = append(errors, ValidationError(fmt.Errorf("Track %d: missing EndOfTrack event", track.TrackNumber)))
		}
	}

	// Program Bank
	for _, track := range smf.Tracks {
		var last any
		for i, e := range track.Events {
			event := e.Event
			c := []*midievent.Controller{nil, nil}

			if cx, ok := last.(midievent.Controller); ok {
				c[0] = &cx
			}

			if cx, ok := event.(midievent.Controller); ok {
				c[1] = &cx
			}

			if c[0] != nil && c[0].Controller.ID == 0x00 && (c[1] == nil || c[1].Controller.ID != 0x20) {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select MSB' event @%d missing LSB (%s)", track.TrackNumber, i, events.Clean(last))))
			}

			if c[1] != nil && c[1].Controller.ID == 0x20 && (c[0] == nil || c[0].Controller.ID != 0x00) {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select LSB' event @%d missing MSB (%s)", track.TrackNumber, i+1, events.Clean(e))))
			}

			if c[0] != nil && c[0].Controller.ID == 0x00 && c[1] != nil && c[1].Controller.ID == 0x20 && c[0].Channel != c[1].Channel {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select MSB' event @%d LSB on another channel (%s)", track.TrackNumber, i, events.Clean(last))))
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: 'Bank Select LSB' event @%d MSB on another channel (%s)", track.TrackNumber, i+1, events.Clean(e))))
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

	if len(smf.Tracks) > 0 {
		track := smf.Tracks[0]
		for _, e := range track.Events {
			if !events.IsTrack0Event(e) {
				errors = append(errors, ValidationError(fmt.Errorf("Track 0: unexpected event (%v)", events.Clean(e))))
			}
		}
	}

	for _, track := range smf.Tracks[1:] {
		for _, e := range track.Events {
			if !events.IsTrack1Event(e) {
				errors = append(errors, ValidationError(fmt.Errorf("Track %d: unexpected event (%s)", track.TrackNumber, events.Clean(e))))
			}
		}
	}

	return errors
}
