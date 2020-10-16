package perftest

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TestConfig struct {
	Factories []string
}

type Config map[string]TestConfig

func (c *Config) Load(filename string) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = jsonFile.Close()
	}()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return json.Unmarshal([]byte(byteValue), &c)
}
