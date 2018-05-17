package models

import (
	"os"
	"testing"

	"github.com/beme/abide"
)

func TestMain(m *testing.M) {
	// Optional: to set a custom directory to load snapshots from,
	// set SnapshotsDir to a path relative to tests before calling
	// any Assert functions
	// here we change the default "__snapshots__" to "testdata"
	abide.SnapshotsDir = "testdata"

	exit := m.Run()
	abide.Cleanup()
	os.Exit(exit)
}

func TestPost(t *testing.T) {
	p := &Post{"Foo", "Bar"}
	abide.Assert(t, "person", p)
}
