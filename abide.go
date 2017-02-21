package abide

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	snapshots SnapshotMap
	args      *arguments
)

func init() {
	// TODO:
	//  - [ ] parse arguments
	//  - [ ] find and load existing snapshots
	args = getArguments()
}

func findExistingSnapshot(id string) (*Snapshot, error) {
	var snapshot *Snapshot

	dir, err := findOrCreateSnapshotDirectory()
	if err != nil {
		return nil, err
	}

	// search inside __snapshots__ dir.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, ErrUnableToReadSnapshotDirectory
	}

	snapshotPaths := []string{}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if filepath.Ext(path) == snapshotExt {
			snapshotPaths = append(snapshotPaths, path)
		}
	}

	if len(snapshotPaths) > 0 {
		snapshotMap, err := findSnapshots(snapshotPaths)
		if err != nil {
			return nil, err
		}
		var ok bool
		snapshot, ok = snapshotMap[id]
		if ok && snapshot != nil {
			return snapshot, nil
		}
	}

	return nil, nil
}

func findOrCreateSnapshotDirectory() (string, error) {
	// get location of tested file.
	testingPath, err := getTestingPath()
	if err != nil {
		return "", ErrUnableToLocateTestPath
	}

	// search for __snapshots__ dir.
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

func createSnapshot(id, val string) (*Snapshot, error) {
	dir, err := findOrCreateSnapshotDirectory()
	if err != nil {
		return nil, err
	}

	pkg, err := getTestingPackage()
	if err != nil {
		return nil, err
	}

	// search inside __snapshots__ dir.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, ErrUnableToReadSnapshotDirectory
	}

	snapshotPath := filepath.Join(dir, fmt.Sprintf("%s%s", pkg, snapshotExt))

	var file *os.File
	var snapshots SnapshotMap
	if len(files) == 0 {
		file, err = os.Create(snapshotPath)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(snapshotPath, os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}

		snapshots, err = findSnapshots([]string{snapshotPath})
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	snapshot := &Snapshot{
		Id:    id,
		Value: val,
	}
	snapshots[id] = snapshot

	err = EncodeSnapshots(file, snapshots)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func findSnapshots(paths []string) (SnapshotMap, error) {
	var snapshotMap = make(map[string]*Snapshot)
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

			snapshots, err := DecodeSnapshots(file)
			if err != nil {
				return
			}

			mutex.Lock()
			for id, snapshot := range snapshots {
				snapshotMap[id] = snapshot
			}
			mutex.Unlock()
		}(paths[i])
	}
	wg.Wait()

	return snapshotMap, nil
}

func getTestingPath() (string, error) {
	return os.Getwd()
}

func getTestingPackage() (string, error) {
	dir, err := getTestingPath()
	if err != nil {
		return "", err
	}

	// TODO:
	//  - parse go test method to figure out if a subdirectory
	//    is being tested.
	return filepath.Base(dir), nil
}
