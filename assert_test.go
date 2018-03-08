package abide

import "testing"

var contentTypeTestCases = []struct {
	input  string
	output bool
}{
	{"application/json", true},
	{"application/json; charset=utf-8", true},
	{"application/vnd.foo.bar.v2+json", true},
	{"application/application/json", false},
	{"application/json/json", false},
	{"application/jsoner; charset=utf-8", false},
	{"application/jsoner", false},
	{"application/vnd.foo.bar.v2+jsoner", false},
	{"application/xml", false},
	{"text/html", false},
	{"", false},
}

func TestContentTypeIsJSON(test *testing.T) {
	for _, testCase := range contentTypeTestCases {

		result := contentTypeIsJSON(testCase.input)

		if result != testCase.output {
			test.Errorf("contentTypeIsJSON(\"%s\" unexpected result. Got=%t, Want=%t", testCase.input, result, testCase.output)
		}
	}
}
