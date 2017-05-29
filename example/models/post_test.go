package models

import (
	"testing"

	"github.com/beme/abide"
)

func TestPost(t *testing.T) {
	p := &Post{"Foo", "Bar"}
	abide.Assert(t, "person", p)
}
