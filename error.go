package abide

import (
	"errors"
)

var (
	ErrUnableToLocateTestPath          = errors.New("unable to locate test path")
	ErrUnableToCreateSnapshotDirectory = errors.New("unable to create snapshot directory")
	ErrUnableToReadSnapshotDirectory   = errors.New("unable to read snapshot directory")
	ErrUnableToLocateSnapshotById      = errors.New("unable to locate snapshot by id")
)
