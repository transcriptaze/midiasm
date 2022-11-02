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
	tag    types.Tag
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

func (e event) Tag() string {
	return fmt.Sprintf("%v", e.tag)
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

// type xxx interface {
// 	*SequenceNumber
// }

// func xyz[T xxx](tick uint64, delta uint32, bytes []byte, f func(uint64, uint32, []byte) (T, error)) (any, error) {
// 	return f(tick, delta, bytes)
// }

// https://stackoverflow.com/questions/71132124/how-to-solve-interface-method-must-have-no-type-parameters
// type pqr func[E xxx](uint64,uint32,[]byte) (E,error)

// var factory2 = map[types.MetaEventType]func(uint64, uint32, []byte) (xxx, error){
// 	0x00: UnmarshalSequenceNumber,
// }

var factory = map[types.MetaEventType]func(uint64, uint32, []byte) (any, error){
	types.TypeSequenceNumber: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalSequenceNumber(tick, delta, bytes)
	},

	types.TypeText: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalText(tick, delta, bytes)
	},

	types.TypeCopyright: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalCopyright(tick, delta, bytes)
	},

	types.TypeTrackName: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTrackName(tick, delta, bytes)
	},

	types.TypeInstrumentName: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalInstrumentName(tick, delta, bytes)
	},

	types.TypeLyric: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalLyric(tick, delta, bytes)
	},

	types.TypeMarker: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMarker(tick, delta, bytes)
	},

	types.TypeCuePoint: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalCuePoint(tick, delta, bytes)
	},

	types.TypeProgramName: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalProgramName(tick, delta, bytes)
	},

	types.TypeDeviceName: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalDeviceName(tick, delta, bytes)
	},

	types.TypeMIDIChannelPrefix: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMIDIChannelPrefix(tick, delta, bytes)
	},

	types.TypeMIDIPort: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMIDIPort(tick, delta, bytes)
	},

	types.TypeEndOfTrack: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalEndOfTrack(tick, delta, bytes)
	},

	types.TypeTempo: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTempo(tick, delta, bytes)
	},

	types.TypeSMPTEOffset: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalSMPTEOffset(tick, delta, bytes)
	},

	types.TypeTimeSignature: func(tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTimeSignature(tick, delta, bytes)
	},
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

	if f, ok := factory[eventType]; ok {
		return f(tick, delta, data)
	}

	switch eventType {
	case 0x59:
		return NewKeySignature(ctx, tick, delta, data)

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
