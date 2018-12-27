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

// Interface is syntactic sugar. It is a helper that converts any type to an Assertable.
func Interface(i interface{}) Assertable {
	// include the type in the string that is asserted to avoid suprises
	return assertableString(fmt.Sprintf("%T %+v", i, i))
}
