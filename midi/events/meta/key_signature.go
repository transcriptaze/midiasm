package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type KeyType uint8

func (k KeyType) String() string {
	if k == 1 {
		return "minor"
	}

	return "major"
}

type KeySignature struct {
	Tag         string
	Status      types.Status
	Type        types.MetaEventType
	Accidentals int8
	KeyType     KeyType
	Key         string
}

var major_keys = map[int8]string{
	7:  `C♯ major`,
	6:  `F♯ major`,
	5:  `B major`,
	4:  `E major`,
	3:  `A major`,
	2:  `D major`,
	1:  `G major`,
	0:  `C major`,
	-1: `F major`,
	-2: `B♭ major`,
	-3: `E♭ major`,
	-4: `A♭ major`,
	-5: `D♭ major`,
	-6: `G♭ major`,
	-7: `C♭ major`,
}

var minor_keys = map[int8]string{
	7:  `A♯ minor`,
	6:  `D♯ minor`,
	5:  `G♯ minor`,
	4:  `C♯ minor`,
	3:  `F♯ minor`,
	2:  `B minor`,
	1:  `E minor`,
	0:  `A minor`,
	-1: `D minor`,
	-2: `G minor`,
	-3: `C minor`,
	-4: `F minor`,
	-5: `B♭ minor`,
	-6: `E♭ minor`,
	-7: `A♭ minor`,
}

var notes = map[int8][]string{
	7:  []string{`C♯`, `D♯`, `E♯`, `F♯`, `G♯`, `A♯`, `B♯`},
	6:  []string{`C♯`, `D♯`, `E♯`, `F♯`, `G♯`, `A♯`, `B`},
	5:  []string{`C♯`, `D♯`, `E`, `F♯`, `G♯`, `A♯`, `B`},
	4:  []string{`C♯`, `D♯`, `E`, `F♯`, `G♯`, `A`, `B`},
	3:  []string{`C♯`, `D`, `E`, `F♯`, `G♯`, `A`, `B`},
	2:  []string{`C♯`, `D`, `E`, `F♯`, `G`, `A`, `B`},
	1:  []string{`C`, `D`, `E`, `F♯`, `G`, `A`, `B`},
	0:  []string{`C`, `D`, `E`, `F`, `G`, `A`, `B`},
	-1: []string{`C`, `D`, `E`, `F`, `G`, `A`, `B♭`},
	-2: []string{`C`, `D`, `E♭`, `F`, `G`, `A`, `B♭`},
	-3: []string{`C`, `D`, `E♭`, `F`, `G`, `A♭`, `B♭`},
	-4: []string{`C`, `D♭`, `E♭`, `F`, `G`, `A♭`, `B♭`},
	-5: []string{`C`, `D♭`, `E♭`, `F`, `G♭`, `A♭`, `B♭`},
	-6: []string{`C♭`, `D♭`, `E♭`, `F`, `G♭`, `A♭`, `B♭`},
	-7: []string{`C♭`, `D♭`, `E♭`, `F♭`, `G♭`, `A♭`, `B♭`},
}

var scales = [][]string{
	{`C`, `C♯`, `D`, `D♯`, `E`, `F`, `F♯`, `G`, `G♯`, `A`, `A♯`, `B`},
	{`C`, `D♭`, `D`, `E♭`, `E`, `F`, `G♭`, `G`, `A♭`, `A`, `B♭`, `B`},
}

func NewKeySignature(ctx *context.Context, bytes []byte) (*KeySignature, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(bytes))
	}

	accidentals := int8(bytes[0])
	keyType := bytes[1]
	if keyType != 0 && keyType != 1 {
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0,1]", keyType)
	}

	key := ""
	switch keyType {
	case 0:
		if signature, ok := major_keys[accidentals]; !ok {
			return nil, fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = signature
		}
	case 1:
		if signature, ok := minor_keys[accidentals]; !ok {
			return nil, fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = signature
		}
	}

	if accidentals < 0 {
		ctx.UseFlats()
	} else {
		ctx.UseSharps()
	}

	return &KeySignature{
		Tag:         "KeySignature",
		Status:      0xff,
		Type:        0x59,
		Accidentals: accidentals,
		KeyType:     KeyType(keyType),
		Key:         key,
	}, nil
}

func (k *KeySignature) Transpose(ctx *context.Context, steps int) {
	// v := int(k.Accidentals) + steps
	// v %= 12

	// // 0 + 8  => 8
	// // 8 % 12 => 8
	// // 8 - 14 => -6

	// k.Accidentals = int8(v)

	// switch k.KeyType {
	// case 0:
	// 	k.Key, _ = major_keys[k.Accidentals]
	// case 1:
	// 	k.Key, _ = minor_keys[k.Accidentals]
	// }

	// if k.Accidentals < 0 {
	// 	ctx.UseFlats()
	// } else {
	// 	ctx.UseSharps()
	// }
}
