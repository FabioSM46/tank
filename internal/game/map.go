package game

import (
	"image/color" // Ensure this is imported for color definitions
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
    TileSize     = 32  // Size of each tile in pixels
    MapWidth     = 25  // Number of horizontal tiles
    MapHeight    = 19  // Number of vertical tiles
    EmptyTile    = 0   // Represents empty space
    WallTile     = 1   // Represents a wall or obstacle
)

type GameMap struct {
    grid [][]int
}

func NewGameMap() *GameMap {
    m := &GameMap{
        grid: make([][]int, MapHeight),
    }

    for i := range m.grid {
        m.grid[i] = make([]int, MapWidth)
    }

    m.Generate()
    return m
}

// Generate randomly places walls on the map
func (m *GameMap) Generate() {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

    for y := 0; y < MapHeight; y++ {
        for x := 0; x < MapWidth; x++ {
            if rnd.Float64() < 0.1 { // 10% chance to place a wall
                m.grid[y][x] = WallTile
            } else {
                m.grid[y][x] = EmptyTile
            }
        }
    }
}

// Draw renders the map to the screen
func (m *GameMap) Draw(screen *ebiten.Image) {
    for y := 0; y < MapHeight; y++ {
        for x := 0; x < MapWidth; x++ {
            tile := m.grid[y][x]
            var col color.Color

            switch tile {
            case EmptyTile:
                col = color.RGBA{0, 128, 0, 255} // Green for empty space
            case WallTile:
                col = color.RGBA{139, 69, 19, 255} // Brown for walls
            }

            // Draw each tile as a rectangle
            ebitenutil.DrawRect(screen, float64(x*TileSize), float64(y*TileSize), TileSize, TileSize, col)
        }
    }
}
