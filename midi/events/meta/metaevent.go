package metaevent

import (
	"bytes"
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/types"
)

type event struct {
	tick   uint64
	delta  uint32
	bytes  []byte
	Tag    string
	Status types.Status
	Type   types.MetaEventType
}

func (e event) Tick() uint64 {
	return e.tick
}

func (e event) Delta() uint32 {
	return e.delta
}

func (e event) Bytes() []byte {
	return e.bytes
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

func Parse(ctx *context.Context, r io.ByteReader, tick uint64, delta uint32) (any, error) {
	status, err := r.ReadByte()
	if err != nil {
		return nil, err
	} else if status != 0xFF {
		return nil, fmt.Errorf("Invalid MetaEvent tag (%v): expected 'FF'", status)
	}

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	eventType := types.MetaEventType(b & 0x7F)

	switch eventType {
	case 0x00:
		return NewSequenceNumber(data)

	case 0x01:
		return NewText(data)

	case 0x02:
		return NewCopyright(data)

	case 0x03:
		return NewTrackName(tick, delta, data), nil

	case 0x04:
		return NewInstrumentName(data)

	case 0x05:
		return NewLyric(data)

	case 0x06:
		return NewMarker(data)

	case 0x07:
		return NewCuePoint(data)

	case 0x08:
		return NewProgramName(data)

	case 0x09:
		return NewDeviceName(data)

	case 0x20:
		return NewMIDIChannelPrefix(data)

	case 0x21:
		return NewMIDIPort(data)

	case 0x51:
		return NewTempo(data)

	case 0x54:
		return NewSMPTEOffset(data)

	case 0x58:
		return NewTimeSignature(data)

	case 0x59:
		return NewKeySignature(ctx, data)

	case 0x2f:
		return NewEndOfTrack(data)

	case 0x7f:
		return NewSequencerSpecificEvent(data)
	}

	return nil, fmt.Errorf("Unrecognised META event: %v", eventType)
}

func concat(list ...[]byte) []byte {
	bytes := []byte{}

	for _, b := range list {
		bytes = append(bytes, b...)
	}

	return bytes
}
