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
	allSnapshots snapshots
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

// Cleanup is an optional method which will execute cleanup operations
// affiliated with abide testing, such as pruning snapshots.
func Cleanup() error {
	for _, s := range allSnapshots {
		if !s.evaluated && args.shouldUpdate && !args.singleRun {
			s.shouldRemove = true
			fmt.Printf("Removing unused snapshot `%s`\n", s.id)
		}
	}

	return allSnapshots.save()
}

// snapshotID represents the unique identifier for a snapshot.
type snapshotID string

// isValid verifies whether the snapshotID is valid. An
// identifier is considered invalid if it is already in use
// or it is malformed.
func (s *snapshotID) isValid() bool {
	return true
}

// snapshot represents the expected value of a test, identified by an id.
type snapshot struct {
	id           snapshotID
	value        string
	path         string
	evaluated    bool
	shouldRemove bool
}

// snapshots represents a map of snapshots by id.
type snapshots map[snapshotID]*snapshot

// save writes all snapshots to their designated files.
func (s snapshots) save() error {
	snapshotsByPath := map[string][]*snapshot{}
	for _, snap := range s {
		_, ok := snapshotsByPath[snap.path]
		if !ok {
			snapshotsByPath[snap.path] = []*snapshot{}
		}
		snapshotsByPath[snap.path] = append(snapshotsByPath[snap.path], snap)
	}

	for path, snaps := range snapshotsByPath {
		if path == "" {
			continue
		}

		snapshotMap := snapshots{}
		for _, snap := range snaps {
			if snap.shouldRemove {
				continue
			}
			snapshotMap[snap.id] = snap
		}
		data, err := encode(snapshotMap)
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

// decode decides a slice of bytes to retrieve a Snapshots object.
func decode(data []byte) (snapshots, error) {
	snaps := make(snapshots)

	snapshotsStr := strings.Split(string(data), snapshotSeparator)
	for _, s := range snapshotsStr {
		if s == "" {
			continue
		}

		components := strings.SplitAfterN(s, "\n", 2)
		id := snapshotID(strings.TrimSuffix(components[0], " */\n"))
		val := strings.TrimSpace(components[1])
		snaps[id] = &snapshot{
			id:    id,
			value: val,
		}
	}

	return snaps, nil
}

// encode encodes a snapshots object into a slice of bytes.
func encode(snaps snapshots) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	ids := []string{}
	for id := range snaps {
		ids = append(ids, string(id))
	}

	sort.Strings(ids)

	data := ""
	for _, id := range ids {
		s := snaps[snapshotID(id)]

		data += fmt.Sprintf("%s%s", snapshotSeparator, string(s.id)) + " */\n"
		data += s.value + "\n\n"
	}

	_, err = buf.WriteString(strings.TrimSpace(data))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// loadSnapshots loads all snapshots in the current directory.
func loadSnapshots() (snapshots, error) {
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

// getSnapshot retrieves a snapshot by id.
func getSnapshot(id snapshotID) *snapshot {
	return allSnapshots[id]
}

// createSnapshot creates or updates a Snapshot.
func createSnapshot(id snapshotID, value string) (*snapshot, error) {
	if !id.isValid() {
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

	s := &snapshot{
		id:    id,
		value: value,
		path:  path,
	}
	allSnapshots[id] = s

	err = allSnapshots.save()
	if err != nil {
		return nil, err
	}

	return s, nil
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

func parseSnapshotsFromPaths(paths []string) (snapshots, error) {
	var snaps = make(snapshots)
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

			s, err := decode(data)
			if err != nil {
				return
			}

			mutex.Lock()
			for id, snap := range s {
				snap.path = p
				snaps[id] = snap
			}
			mutex.Unlock()
		}(paths[i])
	}
	wg.Wait()

	return snaps, nil
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
