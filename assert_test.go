package abide

import (
	"testing"
)

func TestContentTypeIsJSON(test *testing.T) {
	contentTypeTestCases := map[string]bool{
		"application/json":                  true,
		"application/json; charset=utf-8":   true,
		"application/vnd.foo.bar.v2+json":   true,
		"application/application/json":      false,
		"application/json/json":             false,
		"application/jsoner; charset=utf-8": false,
		"application/jsoner":                false,
		"application/vnd.foo.bar.v2+jsoner": false,
		"application/xml":                   false,
		"text/html":                         false,
		"":                                  false,
	}

	for input, expectedOutput := range contentTypeTestCases {
		result := contentTypeIsJSON(input)

		if result != expectedOutput {
			test.Errorf("contentTypeIsJSON(\"%s\" unexpected result. Got=%t, Want=%t", input, result, expectedOutput)
		}
	}
}
