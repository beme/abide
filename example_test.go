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
