package abide

import (
	"fmt"
)

type assertableString string

func (s assertableString) String() string {
	return string(s)
}

// String is syntactic sugar. It is a helper that converts a string to an Assertable
func String(s string) Assertable {
	return assertableString(s)
}

// Struct is syntactic sugar. It is a helper that converts a struct to an Assertable.
func Struct(s interface{}) Assertable {
	return assertableString(fmt.Sprintf("%+v", s))
}
