// Package abide is a testing utility for http response snapshots inspired
// by Facebook's Jest.
//
// It is designed to be used alongside Go's existing testing package
// and enable broader coverage of http APIs. When included in version control
// it can provide a historical log of API and application changes over time.
//
// Snapshot
//
// A snapshot is essentially a lockfile representing an http response.
//  /* snapshot: api endpoint */
//  HTTP/1.1 200 OK
//  Connection: close
//  Content-Type: application/json
//
//  {
//    "foo": "bar"
//  }
//
// In addition to testing `http.Response`, abide provides methods for testing
// `io.Reader` and any object that implements `Assertable`.
//
// Snapshots are saved in a directory named __snapshots__ at the root of the package.
// These files are intended to be saved and included in version control.
//
// Creating a Snapshot
//
// Snapshots are automatically generated during the initial test run. For example
// this will create a snapshot identified by "example" for this http.Response.
//  func TestFunction(t *testing.T) {
//     req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
//     w := httptest.NewRecorder()
//     handler(w, req)
//     res := w.Result()
//     abide.AssertHTTPResponse(t, "example", res)
//  }
//
// Comparing and Updating
//
// In subsequent test runs the existing snapshot is compared to the new results.
// In the event they do not match, the test will fail, and the diff will be printed.
// If the change was intentional, the snapshot can be updated.
//  $ go test -- -u
package abide
