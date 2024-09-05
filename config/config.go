package config

import (
	"os"
)

var WorkingDir string

func init() {
	// Initialize WorkingDir to the current working directory
	dir, err := os.Getwd()
	if err != nil {
		panic("Error getting working directory: " + err.Error())
	}
	WorkingDir = dir
}
