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

type SnapshotId string

func (s *SnapshotId) IsValid() bool {
	return true
}

type Snapshot struct {
	Id    SnapshotId
	Value string

	path string
}

type Snapshots map[SnapshotId]*Snapshot

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
		snapshotMap := Snapshots{}
		for _, snapshot := range snapshots {
			snapshotMap[snapshot.Id] = snapshot
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

func Decode(data []byte) (Snapshots, error) {
	snapshots := make(Snapshots)

	snapshotsStr := strings.Split(string(data), snapshotSeparator)
	for _, s := range snapshotsStr {
		if s == "" {
			continue
		}

		components := strings.SplitAfterN(s, "\n", 2)
		id := SnapshotId(strings.TrimSuffix(components[0], " */\n"))
		val := strings.TrimSpace(components[1])
		snapshots[id] = &Snapshot{
			Id:    id,
			Value: val,
		}
	}

	return snapshots, nil
}

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
		s := snapshots[SnapshotId(id)]

		data += fmt.Sprintf("%s%s", snapshotSeparator, string(s.Id)) + " */\n"
		data += s.Value + "\n\n"
	}

	_, err = buf.WriteString(strings.TrimSpace(data))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

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

func getSnapshot(id SnapshotId) *Snapshot {
	return allSnapshots[id]
}

func createSnapshot(id SnapshotId, value string) (*Snapshot, error) {
	if !id.IsValid() {
		return nil, ErrInvalidSnapshotId
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
		Id:    id,
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
		return "", ErrUnableToLocateTestPath
	}

	dir := filepath.Join(testingPath, snapshotsDir)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return "", ErrUnableToCreateSnapshotDirectory
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
