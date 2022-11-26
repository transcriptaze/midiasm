package lib

import (
	"fmt"
)

type TrackNumber uint

func (t TrackNumber) String() string {
	return fmt.Sprintf("%-2d", t)
}
