package models

import (
	"os"
	"testing"

	"github.com/beme/abide"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	abide.Cleanup()
	os.Exit(exit)
}

func TestPost(t *testing.T) {
	p := &Post{"Foo", "Bar"}
	abide.Assert(t, "person", p)
}
