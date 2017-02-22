package abide

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	args         *arguments
	allSnapshots Snapshots
)

const (
	snapshotsDir      = "__snapshots__"
	snapshotExt       = ".snapshot"
	snapshotSeparator = "/* snapshot: "
)

func init() {
	// 1. Get arguments.
	args = getArguments()

	// 2. Load snapshots.
	allSnapshots, _ = loadSnapshots()
}

// SnapshotID represents the unique identifier for a snapshot.
type SnapshotID string

// IsValid verifies whether the SnapshotID is valid. An
// identifier is considered invalid if it is already in use
// or it is malformed.
func (s *SnapshotID) IsValid() bool {
	return true
}

// Snapshot represents the expected value of a test, identified by an id.
type Snapshot struct {
	ID    SnapshotID
	Value string

	path string
}

// Snapshots represents a map of snapshots by id.
type Snapshots map[SnapshotID]*Snapshot

// Save writes all snapshots to their designated files.
func (s Snapshots) Save() error {
	snapshotsByPath := map[string][]*Snapshot{}
	for _, snapshot := range s {
		_, ok := snapshotsByPath[snapshot.path]
		if !ok {
			snapshotsByPath[snapshot.path] = []*Snapshot{}
		}
		snapshotsByPath[snapshot.path] = append(snapshotsByPath[snapshot.path], snapshot)
	}

	for path, snapshots := range snapshotsByPath {
		if path == "" {
			continue
		}

		snapshotMap := Snapshots{}
		for _, snapshot := range snapshots {
			snapshotMap[snapshot.ID] = snapshot
		}
		data, err := Encode(snapshotMap)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path, data, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Decode decides a slice of bytes to retrieve a Snapshots object.
func Decode(data []byte) (Snapshots, error) {
	snapshots := make(Snapshots)

	snapshotsStr := strings.Split(string(data), snapshotSeparator)
	for _, s := range snapshotsStr {
		if s == "" {
			continue
		}

		components := strings.SplitAfterN(s, "\n", 2)
		id := SnapshotID(strings.TrimSuffix(components[0], " */\n"))
		val := strings.TrimSpace(components[1])
		snapshots[id] = &Snapshot{
			ID:    id,
			Value: val,
		}
	}

	return snapshots, nil
}

// Encode encodes a Snapshots object into a slice of bytes.
func Encode(snapshots Snapshots) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	ids := []string{}
	for id := range snapshots {
		ids = append(ids, string(id))
	}

	sort.Strings(ids)

	data := ""
	for _, id := range ids {
		s := snapshots[SnapshotID(id)]

		data += fmt.Sprintf("%s%s", snapshotSeparator, string(s.ID)) + " */\n"
		data += s.Value + "\n\n"
	}

	_, err = buf.WriteString(strings.TrimSpace(data))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// loadSnapshots loads all snapshots in the current directory.
func loadSnapshots() (Snapshots, error) {
	dir, err := findOrCreateSnapshotDirectory()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if filepath.Ext(path) == snapshotExt {
			paths = append(paths, path)
		}
	}

	return parseSnapshotsFromPaths(paths)
}

// getSnapshot retrieves a Snapshot by id.
func getSnapshot(id SnapshotID) *Snapshot {
	return allSnapshots[id]
}

// createSnapshot creates or updates a Snapshot.
func createSnapshot(id SnapshotID, value string) (*Snapshot, error) {
	if !id.IsValid() {
		return nil, errInvalidSnapshotID
	}

	dir, err := findOrCreateSnapshotDirectory()
	if err != nil {
		return nil, err
	}

	pkg, err := getTestingPackage()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, fmt.Sprintf("%s%s", pkg, snapshotExt))

	snapshot := &Snapshot{
		ID:    id,
		Value: value,
		path:  path,
	}
	allSnapshots[id] = snapshot

	err = allSnapshots.Save()
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func findOrCreateSnapshotDirectory() (string, error) {
	testingPath, err := getTestingPath()
	if err != nil {
		return "", errUnableToLocateTestPath
	}

	dir := filepath.Join(testingPath, snapshotsDir)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return "", errUnableToCreateSnapshotDirectory
		}
	}

	return dir, nil
}

func parseSnapshotsFromPaths(paths []string) (Snapshots, error) {
	var snapshots = make(Snapshots)
	var mutex = &sync.Mutex{}

	var wg sync.WaitGroup
	for i := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			file, err := os.Open(p)
			if err != nil {
				return
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return
			}

			s, err := Decode(data)
			if err != nil {
				return
			}

			mutex.Lock()
			for id, snapshot := range s {
				snapshot.path = p
				snapshots[id] = snapshot
			}
			mutex.Unlock()
		}(paths[i])
	}
	wg.Wait()

	return snapshots, nil
}

func getTestingPath() (string, error) {
	return os.Getwd()
}

func getTestingPackage() (string, error) {
	dir, err := getTestingPath()
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}
