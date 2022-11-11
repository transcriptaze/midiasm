package metaevent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type event struct {
	tick   uint64
	delta  uint32
	bytes  []byte
	tag    lib.Tag
	Status lib.Status
	Type   lib.MetaEventType
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

// type xxx interface {
// 	*SequenceNumber
// }

// func xyz[T xxx](tick uint64, delta uint32, bytes []byte, f func(uint64, uint32, []byte) (T, error)) (any, error) {
// 	return f(tick, delta, bytes)
// }

// https://stackoverflow.com/questions/71132124/how-to-solve-interface-method-must-have-no-type-parameters
// type pqr func[E xxx](uint64,uint32,[]byte) (E,error)

// var factory2 = map[lib.MetaEventType]func(uint64, uint32, []byte) (xxx, error){
// 	0x00: UnmarshalSequenceNumber,
// }

var factory = map[lib.MetaEventType]func(*context.Context, uint64, uint32, []byte) (any, error){
	lib.TypeSequenceNumber: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalSequenceNumber(tick, delta, bytes)
	},

	lib.TypeText: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalText(tick, delta, bytes)
	},

	lib.TypeCopyright: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalCopyright(tick, delta, bytes)
	},

	lib.TypeTrackName: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTrackName(tick, delta, bytes)
	},

	lib.TypeInstrumentName: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalInstrumentName(tick, delta, bytes)
	},

	lib.TypeLyric: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalLyric(tick, delta, bytes)
	},

	lib.TypeMarker: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMarker(tick, delta, bytes)
	},

	lib.TypeCuePoint: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalCuePoint(tick, delta, bytes)
	},

	lib.TypeProgramName: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalProgramName(tick, delta, bytes)
	},

	lib.TypeDeviceName: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalDeviceName(tick, delta, bytes)
	},

	lib.TypeMIDIChannelPrefix: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMIDIChannelPrefix(tick, delta, bytes)
	},

	lib.TypeMIDIPort: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalMIDIPort(tick, delta, bytes)
	},

	lib.TypeEndOfTrack: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalEndOfTrack(tick, delta, bytes)
	},

	lib.TypeTempo: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTempo(tick, delta, bytes)
	},

	lib.TypeSMPTEOffset: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalSMPTEOffset(tick, delta, bytes)
	},

	lib.TypeKeySignature: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		if ks, err := UnmarshalKeySignature(tick, delta, bytes); err != nil {
			return ks, err
		} else {
			if ctx != nil {
				if ks.Accidentals < 0 {
					ctx.UseFlats()
				} else {
					ctx.UseSharps()
				}
			}

			return ks, nil
		}
	},

	lib.TypeTimeSignature: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalTimeSignature(tick, delta, bytes)
	},

	lib.TypeSequencerSpecificEvent: func(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (any, error) {
		return UnmarshalSequencerSpecificEvent(tick, delta, bytes)
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

	eventType := lib.MetaEventType(b & 0x7F)

	if f, ok := factory[eventType]; ok {
		return f(ctx, tick, delta, data)
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
