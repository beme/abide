package abide

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = "abide.json"
)

type config struct {
	HeaderDefaults map[string]interface{} `json:"headerDefaults"`
	Defaults       map[string]interface{} `json:"defaults"`
}

func getConfig() (*config, error) {
	path, err := getTestingPath()
	if err != nil {
		return nil, err
	}

	return parseConfig(filepath.Join(path, configFileName))
}

func parseConfig(path string) (*config, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var c *config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return c, err
}
