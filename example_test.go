package abide_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/beme/abide"
)

var (
	handler = func(*httptest.ResponseRecorder, *http.Request) {}
	t       = &testing.T{}
)

func ExampleAssertHTTPResponse() {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	res := w.Result()
	abide.AssertHTTPResponse(t, "http response", res)
}

func ExampleAssertReader() {
	file, _ := os.Open("/path/to/file")
	abide.AssertReader(t, "io reader", file)
}

func ExampleString() {
	myString := "this is a string I want to snapshot"
	abide.Assert(t, "assertable string", abide.String(myString))
}

func ExampleInterface() {
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
	abide.Assert(t, "assertable struct", abide.Interface(myStruct))
}
