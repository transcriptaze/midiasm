package types

import (
	"fmt"
	"regexp"
	"strconv"
)

type Channel uint8

func (c Channel) String() string {
	return fmt.Sprintf("%d", byte(c))
}

func (c *Channel) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile("channel:([0-9]+)")

	if match := re.FindStringSubmatch(string(bytes)); match == nil || len(match) != 2 {
		return fmt.Errorf("missing channel field")
	} else if channel, err := strconv.ParseUint(match[1], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("channel out of range (%v)", channel)
	} else {
		*c = Channel(channel)
	}

	return nil
}
