package abide

import (
	"os"
	"reflect"
	"testing"
)

func testingCleanup() {
	os.RemoveAll(snapshotsDir)
}

func testingSnapshot(id, value string) *snapshot {
	snapshot, err := createSnapshot(snapshotID(id), value)
	if err != nil {
		panic(err)
	}

	return snapshot
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

	sNew, err := loadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, sNew) {
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