package abide

import (
	"fmt"
)

type assertableString string

func (s assertableString) String() string {
	return string(s)
}

// String is syntactic sugar. It is a helper that converts a string to an AssertableString
func String(s string) assertableString {
	return assertableString(s)
}

// Struct is syntactic sugar. It is a helper that converts a struct to an AssertableString using "%+v" formatting.
func Struct(s interface{}) assertableString {
	return assertableString(fmt.Sprintf("%+v", s))
}
