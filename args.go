package abide

import (
	"os"
	"strings"
)

type arguments struct {
	shouldUpdate bool
	singleRun    bool
}

func getArguments() *arguments {
	args := &arguments{}
	for _, arg := range os.Args {
		argList := strings.Split(arg, "=")
		if len(argList) > 0 {
			arg = argList[0]
		}
		switch arg {
		case "-u":
			args.shouldUpdate = true
			break
		case "-test.run":
			args.singleRun = true
			break
		}
	}

	return args
}
