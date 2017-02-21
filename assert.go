package abide

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func AssertHttpResponse(t *testing.T, id string, w *http.Response) {
	snapshot, err := findExistingSnapshot(id)
	if err != nil {
		t.Fatal(err)
	}

	body, err := httputil.DumpResponse(w, true)
	if err != nil {
		t.Fatal(err)
	}

	if snapshot == nil {
		fmt.Println("Creating new snapshot")
		_, err = createSnapshot(id, string(body))
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	compareResults(t, string(body), snapshot.Value)
}

func compareResults(t *testing.T, new, existing string) {
	dmp := diffmatchpatch.New()
	dmp.PatchMargin = 20
	allDiffs := dmp.DiffMain(new, existing, false)
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
