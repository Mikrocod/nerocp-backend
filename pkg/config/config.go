package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// configuration map
	config map[string]map[string]string
)

// Get configuration value
func Get(parent, child string) string {
	return config[parent][child]
}

// Set configuration value and save
func Set(parent, child, value string) {
	// set value
	config[parent][child] = value

	// marshal
	data, errMarshal := json.Marshal(config)
	shorts.Check(errMarshal)

	// write to file
	errWrite := ioutil.WriteFile("config.json", data, os.ModePerm)
	shorts.Check(errWrite)
}

// LoadConfig load json config
func LoadConfig() {
	// read file
	data, errRead := ioutil.ReadFile("config.json")
	shorts.Check(errRead)

	// unmarshal data to config
	errUnmarshal := json.Unmarshal(data, &config)
	shorts.Check(errUnmarshal)
}
