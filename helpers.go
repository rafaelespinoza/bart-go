package bart

import (
	"strconv"
)

type boolish bool

func (b *boolish) UnmarshalJSON(data []byte) error {
	str := string(data)

	if val, err := strconv.ParseBool(str); err == nil {
		*b = boolish(val)
		return nil
	} else {
		return err
	}
}

func isPresent(target string, data []string) bool {
	return indexOf(target, data) >= 0
}

func indexOf(target string, data []string) int {
	for i, v := range data {
		if v == target {
			return i
		}
	}

	return -1 // not found
}
