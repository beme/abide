package abide

import (
	"io"
	"io/ioutil"
	"strings"
)

const (
	snapshotsDir      = "__snapshots__"
	snapshotExt       = ".snapshot"
	snapshotSeparator = "/* snapshot */"
)

// Snapshot represents the expected value of a test, identified by id.
type Snapshot struct {
	Id    string
	Value string
}

// Snapshots represents a map of snapshots by id.
type SnapshotMap map[string]*Snapshot

// EncodeSnapshots formats and writes a SnapshotMap to an io.Writer.
func EncodeSnapshots(w io.Writer, snapshots SnapshotMap) error {
	var err error
	for _, s := range snapshots {
		_, err = w.Write([]byte(snapshotSeparator + "\n"))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(s.Id + "\n"))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(s.Value + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeSnapshots parses the contents of an io.Reader into a SnapshotMap.
func DecodeSnapshots(r io.Reader) (SnapshotMap, error) {
	var snapshots = make(SnapshotMap)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return snapshots, err
	}

	snapshotsStr := strings.Split(string(data), snapshotSeparator)
	for _, s := range snapshotsStr {
		if s == "" {
			continue
		}
		components := strings.SplitAfterN(s, "\n", 3)
		id := strings.TrimSpace(strings.Trim(components[1], "\n"))
		snapshots[id] = &Snapshot{
			Id:    id,
			Value: strings.TrimSpace(components[2]),
		}
	}

	return snapshots, nil
}
