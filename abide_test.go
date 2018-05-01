package abide

import (
	"os"
	"reflect"
	"testing"
)

func testingCleanup() {
	os.RemoveAll(SnapshotsDir)
}

func testingSnapshot(id, value string) *snapshot {
	snapshot, err := createSnapshot(snapshotID(id), value)
	if err != nil {
		panic(err)
	}

	return snapshot
}

func TestCleanup(t *testing.T) {
	defer testingCleanup()

	_ = testingSnapshot("1", "A")

	// If shouldUpdate = false, the snapshot must remain.
	err := Cleanup()
	if err != nil {
		t.Fatal(err)
	}

	err = LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	snapshot := getSnapshot("1")
	if snapshot == nil {
		t.Fatal("Expected snapshot[1] to exist.")
	}

	// If shouldUpdate = true and singleRun = false, the snapshot must be removed.
	args.shouldUpdate = true
	args.singleRun = false
	err = Cleanup()
	if err != nil {
		t.Fatal(err)
	}

	err = LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	snapshot = getSnapshot("1")
	if snapshot != nil {
		t.Fatal("Expected snapshot[1] to be removed.")
	}
}

func TestSnapshotIDIsValid(t *testing.T) {
	id := snapshotID("1")
	if !id.isValid() {
		t.Fatalf("Expected true, instead got %t.", id.isValid())
	}
}

func TestSnapshotsSave(t *testing.T) {
	defer testingCleanup()

	sA := testingSnapshot("1", "A")
	sB := testingSnapshot("2", "B")

	s := &snapshots{
		"1": sA,
		"2": sB,
	}

	err := s.save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadSnapshots(t *testing.T) {
	defer testingCleanup()

	sA := testingSnapshot("1", "A")
	sB := testingSnapshot("2", "B")

	s := snapshots{
		"1": sA,
		"2": sB,
	}

	err := s.save()
	if err != nil {
		t.Fatal(err)
	}

	err = LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, allSnapshots) {
		t.Fatalf("Failed to load snapshots correctly.")
	}
}

func TestGetSnapshot(t *testing.T) {
	defer testingCleanup()

	snapshot := testingSnapshot("3", "C")
	if !reflect.DeepEqual(snapshot, getSnapshot(snapshot.id)) {
		t.Fatal("Failed to fetch snapshot correctly.")
	}
}
