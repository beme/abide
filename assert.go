package abide

import (
	"encoding/hex"
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

// ContentType expected response body type.
type ContentType string

// ContentType possible values.
const (
	ContentTypeJSON   ContentType = "json"
	ContentTypeText   ContentType = "text"
	ContentTypeBinary ContentType = "binary"
)

// SetContentType define the response type in the options.
func SetContentType(typ ContentType) func(*AssertOptions) {
	return func(options *AssertOptions) {
		options.ContentType = typ
	}
}

// AssertOptions possible options for the assert.
type AssertOptions struct {
	ContentType ContentType
}

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
func AssertHTTPResponse(t *testing.T, id string, w *http.Response, opts ...func(*AssertOptions)) {
	options := &AssertOptions{}
	for _, opt := range opts {
		opt(options)
	}

	body, err := httputil.DumpResponse(w, true)
	if err != nil {
		t.Fatal(err)
	}

	// keep backward compatibility checking for JSON type when the response type
	// wasn't provided
	if options.ContentType == ContentType("") && contentTypeIsJSON(w.Header.Get("Content-Type")) {
		options.ContentType = ContentTypeJSON
	}
	assertHTTP(t, id, body, options.ContentType)
}

// AssertHTTPRequestOut asserts the value of an http.Request.
// Intended for use when testing outgoing client requests
// See https://golang.org/pkg/net/http/httputil/#DumpRequestOut for more
func AssertHTTPRequestOut(t *testing.T, id string, r *http.Request, opts ...func(*AssertOptions)) {
	options := &AssertOptions{}
	for _, opt := range opts {
		opt(options)
	}

	body, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		t.Fatal(err)
	}

	// keep backward compatibility checking for JSON type when the content type
	// wasn't provided
	if options.ContentType == ContentType("") && contentTypeIsJSON(r.Header.Get("Content-Type")) {
		options.ContentType = ContentTypeJSON
	}
	assertHTTP(t, id, body, options.ContentType)
}

// AssertHTTPRequest asserts the value of an http.Request.
// Intended for use when testing incoming client requests
// See https://golang.org/pkg/net/http/httputil/#DumpRequest for more
func AssertHTTPRequest(t *testing.T, id string, r *http.Request, opts ...func(*AssertOptions)) {
	options := &AssertOptions{}
	for _, opt := range opts {
		opt(options)
	}

	body, err := httputil.DumpRequest(r, true)
	if err != nil {
		t.Fatal(err)
	}

	// keep backward compatibility checking for JSON type when the content type
	// wasn't provided
	if options.ContentType == ContentType("") && contentTypeIsJSON(r.Header.Get("Content-Type")) {
		options.ContentType = ContentTypeJSON
	}
	assertHTTP(t, id, body, options.ContentType)
}

func assertHTTP(t *testing.T, id string, body []byte, contentType ContentType) {
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}

	var data string
	if contentType == ContentTypeBinary {
		data = hex.EncodeToString(body)
	} else {
		data = string(body)
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")

	if config != nil {
		// empty line identifies the end of the HTTP header
		for i, line := range lines {
			if line == "" {
				break
			}

			headerItem := strings.Split(line, ":")
			if def, ok := config.Defaults[headerItem[0]]; ok {
				lines[i] = fmt.Sprintf("%s: %s", headerItem[0], def)
			}
		}
	}

	// If the response body is JSON, indent.
	if contentType == ContentTypeJSON {
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
	}

	data = strings.Join(lines, "\n")
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

	return isVendor && isJSON
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
	var err error
	snapshot := getSnapshot(snapshotID(id))

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
			_, err = updateSnapshot(snapshotID(id), data)
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

func didNotMatchMessage(id, diff string) string {
	msg := "\n\n## Existing snapshot does not match results...\n"
	msg += "## \"" + id + "\"\n\n"
	msg += diff
	msg += "\n\n"
	msg += "If this change was intentional, run tests again, $ go test -v -- -u\n"
	return msg
}

func newSnapshotMessage(id, body string) string {
	msg := "\n\n## New snapshot found...\n"
	msg += "## \"" + id + "\"\n\n"
	msg += body
	msg += "\n\n"
	msg += "To save, run tests again, $ go test -v -- -u\n"
	return msg
}
