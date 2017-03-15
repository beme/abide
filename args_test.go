package abide

import (
	"os"
	"testing"
)

func TestGetArguments(t *testing.T) {
	// test without update
	args := getArguments()
	if args.shouldUpdate {
		t.Fatalf("Expected false, instead got %t", args.shouldUpdate)
	}

	// test with update flag
	os.Args = append(os.Args, "-u")
	args = getArguments()
	if !args.shouldUpdate {
		t.Fatalf("Expected true, instead got %t", args.shouldUpdate)
	}
}
