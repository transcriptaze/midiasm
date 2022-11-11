package metaevent

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type TrackName struct {
	event
	Name string
}

func MakeTrackName(tick uint64, delta uint32, name string) TrackName {
	n := lib.VLF(name)
	v, _ := n.MarshalBinary()

	return TrackName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x03}, v...),
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
		},
		Name: name,
	}
}

func UnmarshalTrackName(tick uint64, delta uint32, bytes []byte) (*TrackName, error) {
	name := string(bytes)
	event := MakeTrackName(tick, delta, name)

	return &event, nil
}

func (t TrackName) MarshalBinary() (encoded []byte, err error) {
	var b bytes.Buffer
	var v []byte

	if err = b.WriteByte(byte(t.Status)); err != nil {
		return
	}

	if err = b.WriteByte(byte(t.Type)); err != nil {
		return
	}

	name := lib.VLF(t.Name)
	if v, err = name.MarshalBinary(); err != nil {
		return
	} else if _, err = b.Write(v); err != nil {
		return
	}

	encoded = b.Bytes()

	return
}

func (t *TrackName) UnmarshalText(bytes []byte) error {
	t.tick = 0
	t.delta = 0
	t.bytes = []byte{}
	t.tag = lib.TagTrackName
	t.Status = 0xff
	t.Type = 0x03

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)TrackName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match != nil && len(match) > 2 {
		if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
			return err
		} else {
			t.delta = uint32(delta)
			t.Name = strings.TrimSpace(match[2])
		}
	}

	return nil
}
