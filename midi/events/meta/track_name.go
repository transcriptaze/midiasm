package metaevent

import (
	"bytes"
)

type TrackName struct {
	event
	Name string
}

func NewTrackName(tick uint64, delta uint32, name []byte) *TrackName {
	N, _ := vlq{uint32(len(name))}.MarshalBinary()

	return &TrackName{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: concat([]byte{0x00, 0xff, 0x03}, N, []byte(name)),

			Tag:    "TrackName",
			Status: 0xff,
			Type:   0x03,
		},
		Name: string(name),
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
