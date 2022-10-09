package types

import (
	"fmt"
	"strings"
)

type Hex []byte
type Delta uint32
type Status byte
type MetaEventType byte
type TrackNumber uint
type Channel byte

func (bytes Hex) String() string {
	hex := ""
	for _, b := range bytes {
		hex += fmt.Sprintf("%02X ", b)
	}

	return strings.TrimSpace(hex)
}

func (d Delta) String() string {
	return fmt.Sprintf("%d", d)
}

func (t Status) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t MetaEventType) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t TrackNumber) String() string {
	return fmt.Sprintf("%-2d", t)
}

func (c Channel) String() string {
	return fmt.Sprintf("%d", byte(c))
}
