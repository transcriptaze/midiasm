package midi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type MTrk struct {
	Tag         string           `json:"tag"`
	TrackNumber lib.TrackNumber  `json:"track-number"`
	Length      uint32           `json:"-"`
	Bytes       lib.Hex          `json:"-"`
	Events      []*events.Event  `json:"events"`
	Context     *context.Context `json:"-"`
}

func NewMTrk() (*MTrk, error) {
	mtrk := MTrk{
		Tag:    "MTrk",
		Length: 0,
		Events: []*events.Event{},
		Bytes:  []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x00},
	}

	return &mtrk, nil
}

func (chunk *MTrk) UnmarshalBinary(data []byte) error {
	tag := string(data[0:4])
	if tag != "MTrk" {
		return fmt.Errorf("Invalid MTrk chunk type (%s): expected 'MTrk'", tag)
	}

	length := binary.BigEndian.Uint32(data[4:8])

	eventlist := make([]*events.Event, 0)
	r := bufio.NewReader(bytes.NewReader(data[8:]))
	tick := uint32(0)
	err := error(nil)
	var e *events.Event = nil

	for err == nil {
		e, err = parse(r, tick, chunk.Context)
		if err == nil && e != nil {
			tick += e.Delta()
			eventlist = append(eventlist, e)

			switch k := e.Event.(type) {
			case metaevent.KeySignature:
				if k.Accidentals < 0 {
					chunk.Context.UseFlats()
				} else {
					chunk.Context.UseSharps()
				}

			case midievent.Controller:
				if k.Controller.ID == 0x00 {
					c := uint8(k.Channel)
					v := uint16(k.Value)
					chunk.Context.ProgramBank[c] = (chunk.Context.ProgramBank[c] & 0x003f) | ((v & 0x003f) << 7)
				}

				if k.Controller.ID == 0x20 {
					c := uint8(k.Channel)
					v := uint16(k.Value)
					chunk.Context.ProgramBank[c] = (chunk.Context.ProgramBank[c] & (0x003f << 7)) | (v & 0x003f)
				}

			case midievent.NoteOff:
				e.Event = k.Format(chunk.Context)

			case midievent.NoteOn:
				chunk.Context.PutNoteOn(k.Channel, k.Note.Value)
				e.Event = k.Format(chunk.Context)

			case midievent.ProgramChange:
				c := uint8(k.Channel)
				e.Event = k.SetBank(chunk.Context.ProgramBank[c])
			}
		}
	}

	if err != io.EOF {
		return err
	}

	chunk.Tag = tag
	chunk.Length = length
	chunk.Events = eventlist
	chunk.Bytes = data

	return nil
}

func parse(r *bufio.Reader, tick uint32, ctx *context.Context) (*events.Event, error) {
	rr := reader{
		reader: r,
		buffer: &bytes.Buffer{},
	}

	delta, err := events.VLQ(rr)
	if err != nil {
		return nil, err
	}

	var status byte

	if b, err := rr.peek(); err != nil {
		return nil, err
	} else if b < 0x80 && ctx.RunningStatus < 0x80 {
		return nil, fmt.Errorf("Unrecognised MIDI event: %02X", b&0xF0)
	} else if b < 0x80 {
		status = byte(ctx.RunningStatus)
	} else if status, err = rr.ReadByte(); err != nil {
		return nil, err
	}

	// ... meta event
	if status == 0xff {
		ctx.RunningStatus = 0x00

		if _, err := rr.ReadByte(); err != nil {
			return nil, err
		} else if _, err := events.VLF(rr); err != nil {
			return nil, err
		}

		e, err := metaevent.Parse(uint64(tick)+uint64(delta), rr.Bytes()...)

		return events.NewEvent(e), err
	}

	// ... SysEx event
	if status == 0xf0 || status == 0xf7 {
		ctx.RunningStatus = 0x00

		if _, err := events.VLF(rr); err != nil {
			return nil, err
		} else {
			if status == 0xf0 && ctx.Casio {
				return nil, fmt.Errorf("Invalid SysExMessage event data: F0 start byte without terminating F7")
			}

			if e, err := sysex.Parse(ctx, uint64(tick)+uint64(delta), rr.Bytes()...); err != nil {
				return nil, err
			} else {
				bytes := rr.Bytes()
				if status == 0xf0 {
					ctx.Casio = bytes[len(bytes)-1] != 0xf7
				} else if status == 0xf7 && ctx.Casio && bytes[len(bytes)-1] == 0xf7 {
					ctx.Casio = false
				}

				return events.NewEvent(e), err
			}
		}
	}

	// ... MIDI event
	var length = map[byte]int{
		0x80: 2,
		0x90: 2,
		0xA0: 1,
		0xB0: 2,
		0xC0: 1,
		0xD0: 1,
		0xE0: 2,
	}

	for i := 0; i < length[status&0xf0]; i++ {
		if _, err := rr.ReadByte(); err != nil {
			return nil, err
		}
	}

	e, err := midievent.Parse(uint64(tick)+uint64(delta), ctx.RunningStatus, rr.Bytes()...)

	ctx.RunningStatus = status

	return events.NewEvent(e), err
}
