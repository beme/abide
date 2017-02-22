package abide

import (
	"os"
)

type arguments struct {
	shouldUpdate bool
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
