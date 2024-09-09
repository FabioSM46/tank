package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/toolchain/src/math/rand"
)

const (
	TileSize  = 32 // Size of each tile in pixels
	MapWidth  = 25 // Number of horizontal tiles
	MapHeight = 19 // Number of vertical tiles
	EmptyTile = 0  // Represents empty space
	WallTile  = 1  // Represents a wall or obstacle
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

func (m *GameMap) CheckCollisionWithWall(px1, py1, px2, py2 float64) bool {
	tileSize := float64(TileSize)
	collisionDetected := false
	// Check if the player's bounding box is outside the map boundaries
	if px1 < 0 || py1 < 0 || px2 > float64(MapWidth)*tileSize || py2 > float64(MapHeight)*tileSize {
		return true // Collision detected with the map boundary
	}
	// Check collision with wall tiles
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			if m.grid[y][x] == WallTile {
				tx1, ty1 := float64(x)*tileSize, float64(y)*tileSize
				tx2, ty2 := tx1+tileSize, ty1+tileSize
				// Check if the player's bounding box intersects with the wall tile
				if px1 < tx2 && px2 > tx1 && py1 < ty2 && py2 > ty1 {
					collisionDetected = true
					break
				}
			}
		}
	}

	return collisionDetected
}

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

func (gm *GameMap) CheckLineOfSight(x1, y1, x2, y2 float64) bool {
	// Calculate the number of steps needed for the line trace
	steps := int(math.Max(math.Abs(x2-x1), math.Abs(y2-y1)) / float64(TileSize))
	if steps == 0 {
		return true // If no steps, consider it a clear line of sight
	}
	dx := (x2 - x1) / float64(steps)
	dy := (y2 - y1) / float64(steps)

	// Trace the line from (x1, y1) to (x2, y2) in small increments
	for i := 0; i <= steps; i++ {
		currentX := x1 + dx*float64(i)
		currentY := y1 + dy*float64(i)

		// Check if the current position is within a wall tile
		if gm.IsWallTileAt(currentX, currentY) {
			// Add debug output to help trace issues
			return false // There's an obstacle in the way
		}
	}

	return true // No obstacle detected
}

func (gm *GameMap) IsWallTileAt(x, y float64) bool {
	tileX := int(x / float64(TileSize))
	tileY := int(y / float64(TileSize))

	// Check for bounds and if the tile is a wall
	if tileX >= 0 && tileX < MapWidth && tileY >= 0 && tileY < MapHeight {
		return gm.grid[tileY][tileX] == WallTile
	}
	return false
}
