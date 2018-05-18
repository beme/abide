# abide

A testing utility for http response snapshots. Inspired by [Jest](https://github.com/facebook/jest).

[![Build Status](https://travis-ci.org/beme/abide.png)](https://travis-ci.org/beme/abide)
[![GoDoc](https://godoc.org/github.com/beme/abide?status.svg)](https://godoc.org/github.com/beme/abide)

## Usage

1. Include abide in your project.

```go
import "github.com/beme/abide"
```

2. Within your test function, capture the response to an http request, set a unique identifier, and assert.

```go
func TestFunction(t *testing.T) {
  req := httptest.NewRequest("GET", "http://example.com/", nil)
  w := httptest.NewRecorder()
  exampleHandler(w, req)
  res := w.Result()
  abide.AssertHTTPResponse(t, "example route", res)
}
```

3. Run your tests.
```shell
$ go test -v
```

4. If the output of your http response does not equal the existing snapshot, the difference will be printed in the test output. If this change was intentional, the snapshot can be updated by including the `-u` flag.
```shell
$ go test -v -- -u
```

Any snapshots created/updated will be located in `package/__snapshots__`.

5. Cleanup

To ensure only the snapshots in-use are included, add the following to `TestMain`. If your application does not have one yet, you can read about `TestMain` usage [here](https://golang.org/pkg/testing/#hdr-Main).

```go
func TestMain(m *testing.M) {
  exit := m.Run()
  abide.Cleanup()
  os.Exit(exit)
}
```

Once included, if the update `-u` flag is used when running tests, any snapshot that is no longer in use will be removed. Note: if a single test is run, pruning _will not occur_.

## Snapshots

A snapshot is essentially a lock file for an http response. Instead of having to manually compare every aspect of an http response to it's expected value, it can be automatically generated and used for matching in subsequent testing.

Here's an example snapshot:

```
/* snapshot: example route */
HTTP/1.1 200 OK
Connection: close
Content-Type: application/json

{
  "key": "value"
}
```

`abide` also supports testing outside of http responses, by providing an `Assert(*testing.T, string, Assertable)` method which will create snapshots for any type that implements `String() string`.

## Example

See `/example` for the usage of `abide` in a basic web server. To run tests, simply `$ go test -v`

## Config

In some cases, attributes in a JSON response can by dynamic (e.g unique id's, dates, etc.), which can disrupt snapshot testing. To resolve this, an `abide.json` file config can be included to override values with defaults. Consider the config in the supplied example project:

```json
{
  "defaults": {
    "updated_at": 0,
    "foo": "foobar"
  }
}
```

When used with `AssertHTTPResponse`, for any response with `Content-Type: application/json`, the key-value pairs in `defaults` will be used to override the JSON response, allowing for consistent snapshot testing.


## Using custom `__snapshot__` directory

To write snapshots to a directory other than the default `__snapshot__`, adjust `abide.SnapshotDir` before any call to an Assert function. See `example/models` package for a working example

```golang
func init() {
  abide.SnapshotDir = "testdata"
}
``` 