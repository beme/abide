package abide

import (
	"os"
)

type arguments struct {
	// ShouldUpdate represents whether the snapshot should
	// be updated if there is a diff.
	ShouldUpdate bool
}

func getArguments() *arguments {
	shouldUpdate := false
	for _, arg := range os.Args {
		if arg == "-u" {
			shouldUpdate = true
			break
		}
	}

	return &arguments{shouldUpdate}
}
