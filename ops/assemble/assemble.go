package assemble

import (
	"io"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
)

type Assembler interface {
	Assemble(io.Reader) ([]byte, error)
}

func fixups(mtrk *midi.MTrk) (*midi.MTrk, error) {
	for ix := range mtrk.Events[1:] {
		previous := mtrk.Events[ix]
		event := mtrk.Events[ix+1]

		if message, ok := previous.Event.(*sysex.SysExMessage); ok {
			if _, ok := event.Event.(*sysex.SysExContinuationMessage); !ok {
				message.Single = true
			}
		}
	}

	for ix := range mtrk.Events[1:] {
		previous := mtrk.Events[ix]
		event := mtrk.Events[ix+1]

		if message, ok := previous.Event.(*sysex.SysExContinuationMessage); ok {
			if _, ok := event.Event.(*sysex.SysExContinuationMessage); !ok {
				message.End = true
			}
		}
	}

	return mtrk, nil
}
