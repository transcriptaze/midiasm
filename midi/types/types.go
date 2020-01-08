package types

import (
	"fmt"
	"strings"
)

type Hex []byte
type Tick uint64
type Delta uint32
type MetaEventType byte
type TrackNumber uint

func (bytes Hex) String() string {
	hex := ""
	for _, b := range bytes {
		hex += fmt.Sprintf("%02X ", b)
	}

	return strings.TrimSpace(hex)
}

func (t Tick) String() string {
	return fmt.Sprintf("%d", t)
}

func (d Delta) String() string {
	return fmt.Sprintf("%d", d)
}

func (t MetaEventType) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t TrackNumber) String() string {
	return fmt.Sprintf("%-2d", t)
}
