package types

import (
	"fmt"
	"strings"
)

type Hex []byte
type Status byte
type MetaEventType byte
type TrackNumber uint

func (bytes Hex) String() string {
	hex := ""
	for _, b := range bytes {
		hex += fmt.Sprintf("%02X ", b)
	}

	return strings.TrimSpace(hex)
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
