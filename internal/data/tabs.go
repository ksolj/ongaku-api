package data

import (
	"fmt"
	"strconv"
)

type Tab string // implementation may change in the future

func (r Tab) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("Tabs: %s", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
