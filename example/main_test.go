package main

import (
	"net/http/httptest"
	"testing"

	"github.com/beme/abide"
)

func TestFunction(t *testing.T) {
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
}
