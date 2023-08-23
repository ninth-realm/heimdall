package http

import (
	"encoding/json"
	"net/mail"
	"reflect"
	"strings"
)

type optional[T any] struct {
	Value   T
	Present bool
}

func (o *optional[T]) UnmarshalJSON(b []byte) error {
	var inner T
	err := json.Unmarshal(b, &inner)
	if err != nil {
		return err
	}

	*o = optional[T]{
		Value:   inner,
		Present: true,
	}

	return nil
}

type emailString string

func (s *emailString) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)

	if err != nil {
		return err
	}

	_, err = mail.ParseAddress(str)
	if err != nil {
		return &json.UnmarshalTypeError{Value: string(b), Type: reflect.TypeOf(*s)}
	}

	*s = emailString(str)

	return nil
}

// nonEmptyString represents a string that contains at least one non-whitespace
// character.
type nonEmptyString string

func (s nonEmptyString) toString() string { return string(s) }

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
