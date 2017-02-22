package abide

import (
	"errors"
)

var (
	errUnableToLocateTestPath          = errors.New("unable to locate test path")
	errUnableToCreateSnapshotDirectory = errors.New("unable to create snapshot directory")
	errUnableToReadSnapshotDirectory   = errors.New("unable to read snapshot directory")
	errUnableToLocateSnapshotByID      = errors.New("unable to locate snapshot by id")
	errInvalidSnapshotID               = errors.New("invalid snapshot id")
)
