package midi

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events/meta"
	"strings"
)

type SMF struct {
	File   string
	MThd   *MThd
	Tracks []*MTrk
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

	return errors
}
