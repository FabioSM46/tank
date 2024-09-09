package utils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/FabioSM46/tank/config"
)

func LoadWindowSettings() WindowSettings {
	file, err := os.Open(filepath.Join(config.WorkingDir, "pkg/json/window_settings.json"))
	if err != nil {
		return WindowSettings{X: 100, Y: 100}
	}
	defer file.Close()

	settings := WindowSettings{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&settings); err != nil {
		return WindowSettings{X: 100, Y: 100}
	}

	return settings
}

func SaveWindowSettings(settings WindowSettings) {
	file, _ := os.Create(filepath.Join(config.WorkingDir, "pkg/json/window_settings.json"))
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(settings)
}
