package abide

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/beme/abide/internal"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Assertable represents an object that can be asserted.
type Assertable interface {
	String() string
}

// Assert asserts the value of an object with implements Assertable.
func Assert(t *testing.T, id string, a Assertable) {
	data := a.String()
	createOrUpdateSnapshot(t, id, data)
}

// AssertHTTPResponse asserts the value of an http.Response.
func AssertHTTPResponse(t *testing.T, id string, w *http.Response) {
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}

	body, err := httputil.DumpResponse(w, true)
	if err != nil {
		t.Fatal(err)
	}

	data := string(body)

	contentType := w.Header.Get("Content-Type")

	// If the response body is JSON, indent.
	if contentTypeIsJSON(contentType) {
		lines := strings.Split(data, "\n")
		jsonStr := lines[len(lines)-1]

		var jsonIface map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &jsonIface)
		if err != nil {
			t.Fatal(err)
		}

		// Clean/update json based on config.
		if config != nil {
			for k, v := range config.Defaults {
				jsonIface = internal.UpdateKeyValuesInMap(k, v, jsonIface)
			}
		}

		out, err := json.MarshalIndent(jsonIface, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		lines[len(lines)-1] = string(out)
		data = strings.Join(lines, "\n")
	}

	createOrUpdateSnapshot(t, id, data)
}

func contentTypeIsJSON(contentType string) bool {
	contentTypeParts := strings.Split(contentType, ";")
	firstPart := contentTypeParts[0]

	isPlainJSON := firstPart == "application/json"
	if isPlainJSON {
		return isPlainJSON
	}

	isVendor := strings.HasPrefix(firstPart, "application/vnd.")

	isJSON := strings.HasSuffix(firstPart, "+json")

	if isVendor && isJSON {
		return true
	}

	return false
}

// AssertReader asserts the value of an io.Reader.
func AssertReader(t *testing.T, id string, r io.Reader) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	createOrUpdateSnapshot(t, id, string(data))
}

func createOrUpdateSnapshot(t *testing.T, id, data string) {
	snapshot := getSnapshot(snapshotID(id))

	var err error
	if snapshot == nil {
		if !args.shouldUpdate {
			t.Error(newSnapshotMessage(id, data))
			return
		}

		fmt.Printf("Creating snapshot `%s`\n", id)
		snapshot, err = createSnapshot(snapshotID(id), data)
		if err != nil {
			t.Fatal(err)
		}
		snapshot.evaluated = true
		return
	}

	snapshot.evaluated = true
	diff := compareResults(t, snapshot.value, strings.TrimSpace(data))
	if diff != "" {
		if snapshot != nil && args.shouldUpdate {
			fmt.Printf("Updating snapshot `%s`\n", id)
			_, err = createSnapshot(snapshotID(id), data)
			if err != nil {
				t.Fatal(err)
			}
			return
		}

		t.Error(didNotMatchMessage(id, diff))
		return
	}
}

func compareResults(t *testing.T, existing, new string) string {
	dmp := diffmatchpatch.New()
	dmp.PatchMargin = 20
	allDiffs := dmp.DiffMain(existing, new, false)
	nonEqualDiffs := []diffmatchpatch.Diff{}
	for _, diff := range allDiffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			nonEqualDiffs = append(nonEqualDiffs, diff)
		}
	}

	if len(nonEqualDiffs) == 0 {
		return ""
	}

	return dmp.DiffPrettyText(allDiffs)
}

func didNotMatchMessage(id string, diff string) string {
	msg := "\n\n## Existing snapshot does not match results...\n"
	msg += "## \"" + id + "\"\n\n"
	msg += diff
	msg += "\n\n"
	msg += "If this change was intentional, run tests again, $ go test -v -- -u\n"
	return msg
}

func newSnapshotMessage(id string, body string) string {
	msg := "\n\n## New snapshot found...\n"
	msg += "## \"" + id + "\"\n\n"
	msg += body
	msg += "\n\n"
	msg += "To save, run tests again, $ go test -v -- -u\n"
	return msg
}
