package lib

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

type Channel uint8

func ParseChannel(s string) (Channel, error) {
	if channel, err := strconv.ParseUint(s, 10, 8); err != nil {
		return 0, err
	} else if channel > 15 {
		return 0, fmt.Errorf("invalid channel (%v)", channel)
	} else {
		return Channel(channel), nil
	}
}

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

func (c Channel) MarshalJSON() ([]byte, error) {
	channel := uint(c)

	return json.Marshal(channel)
}

func (c *Channel) UnmarshalJSON(bytes []byte) error {
	var v uint

	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	} else if v > 15 {
		return fmt.Errorf("Invalid channel (%v)", v)
	}

	*c = Channel(v)

	return nil
}
