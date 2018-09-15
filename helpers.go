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
