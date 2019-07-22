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

func testingSnapshots(count int) snapshots {
	s := make(snapshots, count)
	for i := 0; i < count; i++ {
		id := string(i)
		s[snapshotID(id)] = testingSnapshot(id, id)
	}
	return s
}

func TestCleanup(t *testing.T) {
	defer testingCleanup()

	_ = testingSnapshot("1", "A")

	// If shouldUpdate = false, the snapshot must remain.
	err := Cleanup()
	if err != nil {
		t.Fatal(err)
	}

	err = loadSnapshots()
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

	// call private reloadSnapshots to repeat once-executing function
	err = reloadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	snapshot = getSnapshot("1")
	if snapshot != nil {
		t.Fatal("Expected snapshot[1] to be removed.")
	}
}

func TestCleanupUpdate(t *testing.T) {
	defer testingCleanup()

	_ = testingSnapshot("1", "A")
	t2 := &testing.T{}
	createOrUpdateSnapshot(t2, "1", "B")

	snapshot := getSnapshot("1")
	if snapshot == nil {
		t.Fatal("Expected snapshot[1] to exist.")
	}

	// If shouldUpdate = true and singleRun = false, the snapshot must be removed.
	args.shouldUpdate = true
	args.singleRun = false
	err := Cleanup()
	if err != nil {
		t.Fatal(err)
	}

	// call private reloadSnapshots to repeat once-executing function
	err = reloadSnapshots()
	if err != nil {
		t.Fatal(err)
	}

	snapshot = getSnapshot("1")
	if snapshot == nil {
		t.Fatal("Expected snapshot[1] to exist.")
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

	err = loadSnapshots()
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

func benchmarkEncode(count int, b *testing.B) {
	defer testingCleanup()
	s := testingSnapshots(count)
	for i := 0; i < b.N; i++ {
		encode(s)
	}
}

func BenchmarkEncode10(b *testing.B)   { benchmarkEncode(10, b) }
func BenchmarkEncode100(b *testing.B)  { benchmarkEncode(100, b) }
func BenchmarkEncode1000(b *testing.B) { benchmarkEncode(1000, b) }
