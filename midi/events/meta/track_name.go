package metaevent

import (
	"bytes"
	"regexp"
	"strings"
)

type TrackName struct {
	event
	Name string
}

func NewTrackName(tick uint64, delta uint32, bytes []byte) *TrackName {
	N, _ := vlq{uint32(len(bytes))}.MarshalBinary()

	return &TrackName{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: concat([]byte{0x00, 0xff, 0x03}, N, []byte(bytes)),

			Tag:    "TrackName",
			Status: 0xff,
			Type:   0x03,
		},
		Name: string(bytes),
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

func (e *TrackName) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.Tag = "TrackName"
	e.Type = 0x03

	re := regexp.MustCompile(`(?i)TrackName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match != nil && len(match) > 1 {
		e.Name = strings.TrimSpace(match[1])
	}

	return nil
}
