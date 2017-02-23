package abide

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// AssertHTTPResponse asserts the value of an http.Response.
func AssertHTTPResponse(t *testing.T, id string, w *http.Response) {
	snapshot := getSnapshot(snapshotID(id))

	body, err := httputil.DumpResponse(w, true)
	if err != nil {
		t.Fatal(err)
	}

	if snapshot == nil {
		fmt.Printf("Creating snapshot `%s`\n", id)
		_, err = createSnapshot(snapshotID(id), string(body))
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	if snapshot != nil && args.shouldUpdate {
		fmt.Printf("Updating snapshot `%s`\n", id)
		_, err = createSnapshot(snapshotID(id), string(body))
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	compareResults(t, snapshot.value, string(body))
}

func compareResults(t *testing.T, existing, new string) {
	dmp := diffmatchpatch.New()
	dmp.PatchMargin = 20
	allDiffs := dmp.DiffMain(existing, new, false)
	nonEqualDiffs := []diffmatchpatch.Diff{}
	for _, diff := range allDiffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			nonEqualDiffs = append(nonEqualDiffs, diff)
		}
	}
	if len(nonEqualDiffs) > 0 {
		msg := "\n\nExisting snapshot does not match results...\n\n"
		msg += dmp.DiffPrettyText(allDiffs)
		msg += "\n\n"
		msg += "If this change was intentional, run tests again, $ go test -v -- -u\n"

		t.Error(msg)
	}
}
