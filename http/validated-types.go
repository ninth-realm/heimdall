package http

import (
	"encoding/json"
	"reflect"
	"strings"
)

// nonEmptyString represents a string that contains at least one non-whitespace
// character.
type nonEmptyString string

func (s *nonEmptyString) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)

	if err != nil {
		return err
	}

	if len(strings.TrimSpace(str)) == 0 {
		return &json.UnmarshalTypeError{Value: string(b), Type: reflect.TypeOf(*s)}
	}

	*s = nonEmptyString(str)

	return nil
}
