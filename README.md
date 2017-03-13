# abide

A testing utility for http response snapshots. Inspired by [jest](https://github.com/facebook/jest).

## Usage

1. Include abide in your project.

```
$ go get github.com/beme/abide
```

2. Within your test function, capture the response to an http request, set a unique identifier, and assert.

```go
func TestFunction(t *testing.T) {
  req := httptest.NewRequest("GET", "http://example.com/", nil)
  w := httptest.NewRecorder()
  firstHandler(w, req)
  res := w.Result()
  abide.AssertHTTPResponse(t, "example route", res)
}
```

3. Run your tests.
```
$ go test -v
```

4. If the output of your http response does not equal the existing snapshot, and this was intentional, the snapshot can be updated.
```
$ go test -v -- -u
```

## Snapshots

A snapshot is essentially a lock file for an http response. Instead of having to manually compare every aspect of an http response to it's expected value, it can automatically generated and used for matching in subsequent testing.

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

Snapshots are located within `project/__snapshots__`.

## Example

See `/example` for the usage of `abide` in a basic web server. To run tests, simply `$ go test -v`

## Config

In some cases, attributes in a JSON response can by dynamic (e.g unique id's, dates, etc.), which can disrupt snapshot testing. To resolve this, an `abide.json` file config file can be included to override values with defaults. Consider the config in the supplied example project:

```json
{
  "defaults": {
    "updated_at": 0,
    "foo": "foobar"
  }
}
```

For any JSON response, the key-value pairs in `defaults` will be used to override, allowing for consistent snapshot testing.
