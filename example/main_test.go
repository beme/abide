package main

import (
	"net/http/httptest"
	"os"
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
