package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/toolchain/src/math/rand"
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

// CheckCollisionWithWall checks if any part of the player's bounding box intersects with a wall tile
func (m *GameMap) CheckCollisionWithWall(px1, py1, px2, py2 float64) bool {
    tileSize := float64(TileSize)
    collisionDetected := false

    for y := 0; y < MapHeight; y++ {
        for x := 0; x < MapWidth; x++ {
            if m.grid[y][x] == WallTile {
                tx1, ty1 := float64(x)*tileSize, float64(y)*tileSize
                tx2, ty2 := tx1+tileSize, ty1+tileSize
                // Check if the player's bounding box intersects with the wall tile
                if px1 < tx2 && px2 > tx1 && py1 < ty2 && py2 > ty1 {
                    collisionDetected = true
                    // Additional code to handle position adjustment
                }
            }
        }
    }
    return collisionDetected
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
            vector.DrawFilledRect(screen, float32(x*TileSize), float32(y*TileSize), float32(TileSize), float32(TileSize), col, false)
        }
    }
}

/* // IsWall checks if a position collides with a wall.
func (m *GameMap) IsWall(x, y float64) bool {
	tileX := int(x) / TileSize 
	tileY := int(y) / TileSize 

	// Check if the position is outside the bounds of the map
	if tileX < 0 || tileY < 0 || tileX >= MapWidth || tileY >= MapHeight {
		return true // Treat out-of-bounds as a wall
	}

	return m.grid[tileY][tileX] == WallTile
} */