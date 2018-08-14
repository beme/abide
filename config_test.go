package abide

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestingConfig() {
	data := []byte(`{
    "headerDefaults": {
      "baz": "quz"
    },
    "defaults": {
      "foo": "bar"
    }
  }`)
	err := ioutil.WriteFile(configFileName, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestGetConfig(t *testing.T) {
	defer testingCleanup()

	// test no config
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}

	TestingConfig()
	defer os.Remove(configFileName)

	// test with config
	config, err = getConfig()
	if err != nil {
		t.Fatal(err)
	}

	if config.HeaderDefaults["baz"] != "quz" {
		t.Fatalf("Expected to find header default value quz, instead got %s.", config.HeaderDefaults["baz"])
	}

	if config.Defaults["foo"] != "bar" {
		t.Fatalf("Expected to find default value bar, instead got %s.", config.Defaults["foo"])
	}
}
