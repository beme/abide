package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/beme/abide"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	abide.Cleanup()
	os.Exit(exit)
}

func TestRequests(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	firstHandler(w, req)
	res := w.Result()
	abide.AssertHTTPResponse(t, "first route", res)

	req = httptest.NewRequest("GET", "http://example.com/", nil)
	w = httptest.NewRecorder()
	secondHandler(w, req)
	res = w.Result()
	abide.AssertHTTPResponse(t, "second route", res)

	req = httptest.NewRequest("GET", "http://example.com/", nil)
	w = httptest.NewRecorder()
	thirdHandler(w, req)
	res = w.Result()
	abide.AssertHTTPResponse(t, "third route", res)
}

func TestReader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	fourthHandler(w, req)
	res := w.Result()
	abide.AssertReader(t, "reader", res.Body)
}

func TestAssertHTTPRequestOut(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", strings.NewReader(`{"message": "expected message"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Expected-Header", "expected header value")

	abide.AssertHTTPRequestOut(t, "http client request", req)
}

func TestAssertHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = "" // httptest servers are spawned on random ports, prevent that from being in the snapshot.
		abide.AssertHTTPRequest(t, "http server request", r)
	}))

	http.Post(server.URL, "application/json", strings.NewReader(`{"message": "expected message"}`))
}

func TestAssertableString(t *testing.T) {
	abide.Assert(t, "assertable string", abide.String("string to be asserted"))
}

func TestAssertableInterface(t *testing.T) {
	type MyStruct struct {
		Field1 string
		Field2 int64
		Field3 bool
		field4 string
	}
	myStruct := MyStruct{
		"String1",
		1234567,
		true,
		"string4",
	}
	abide.Assert(t, "assertable interface", abide.Interface(myStruct))
}
