package game

import (
	"encoding/json"
	"os"
)

// Map represents the game map.
type Map struct {
	Tiles       []string `json:"tiles"`
	PlayerStart struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"playerStart"`
}

// LoadMap loads a map from a JSON file.
func LoadMap(filename string) (*Map, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var m Map
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// GetTile returns the tile at the specified coordinates.
func (m *Map) GetTile(x, y int) rune {
	if y < 0 || y >= len(m.Tiles) || x < 0 || x >= len(m.Tiles[y]) {
		return ' '
	}
	return rune(m.Tiles[y][x])
}
