package metaevent

import (
	"bytes"

	"github.com/transcriptaze/midiasm/midi/types"
)

type TrackName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

type vlq struct {
	v uint32
}

func (v vlq) MarshalBinary() ([]byte, error) {
	buffer := []byte{0, 0, 0, 0, 0}
	b := v.v

	for i := 4; i > 0; i-- {
		buffer[i] = byte(b & 0x7f)
		if b >>= 7; b == 0 {
			return buffer[i:], nil
		}
	}

	buffer[1] |= 0x80
	buffer[0] = byte(b & 0x7f)

	return buffer, nil
}

type vlf struct {
	v []byte
}

func (v vlf) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer
	var u []byte

	N := vlq{uint32(len(v.v))}
	if u, err = N.MarshalBinary(); err != nil {
		return
	} else if _, err = b.Write(u); err != nil {
		return
	}

	if _, err = b.Write(v.v); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}

func NewTrackName(bytes []byte) (*TrackName, error) {
	return &TrackName{
		Tag:    "TrackName",
		Status: 0xff,
		Type:   0x03,
		Name:   string(bytes),
	}, nil
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
