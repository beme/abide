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
