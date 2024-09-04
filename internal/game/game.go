package game

import (
	"image/color"
	"log"
	"os"

	"github.com/FabioSM46/tank/internal/player"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Game represents the current game state.
type Game struct {
	Player *player.Player
	Map    *Map
}

// NewGame initializes and returns a new Game instance.
func NewGame(screenWidth, screenHeight int) *Game {
	g := &Game{}
	
	// Load the map
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	m, err := LoadMap(wd + "/internal/game/maps/map.json")
	if err != nil {
		log.Fatalf("Failed to load map: %v", err)
	}
	g.Map = m
	
	// Initialize player
	playerX := float64(m.PlayerStart.X * utils.PlayerWidth)
	playerY := float64(m.PlayerStart.Y * utils.PlayerHeight)
	g.Player = player.NewPlayer(playerX, playerY)
	
	return g
}

// Update updates the game logic.
func (g *Game) Update() error {
	g.Player.Update()
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the map
	for y, row := range g.Map.Tiles {
		for x, tile := range row {
			switch tile {
			case '#':
				color := color.RGBA{100, 100, 100, 255} // Gray for walls
				vector.DrawFilledRect(screen, float32(x*utils.WallWidth), float32(y*utils.WallHeight), utils.WallWidth, utils.WallHeight, color, false)
			case '.':
				// Empty space, no need to draw anything
			}
		}
	}

	// Draw the player
	g.Player.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
