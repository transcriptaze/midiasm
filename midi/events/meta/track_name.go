package metaevent

import (
	"bytes"

	"github.com/transcriptaze/midiasm/midi/types"
)

type TrackName struct {
	event
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewTrackName(tick uint64, delta uint32, name string) *TrackName {
	N, _ := vlq{uint32(len(name))}.MarshalBinary()
	data := []byte(name)
	bytes := append(append([]byte{0x00, 0xff, 0x03}, N...), data...)

	return &TrackName{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: bytes,
		},
		Tag:    "TrackName",
		Status: 0xff,
		Type:   0x03,
		Name:   name,
	}
}

func (e TrackName) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer
	var v []byte

	if err = b.WriteByte(byte(e.Status)); err != nil {
		return
	}

	if err = b.WriteByte(byte(e.Type)); err != nil {
		return
	}

	name := vlf{[]byte(e.Name)}
	if v, err = name.MarshalBinary(); err != nil {
		return
	} else if _, err = b.Write(v); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}
