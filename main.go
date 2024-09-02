package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game represents the current game state.
type Game struct {
	playerX, playerY float64
}

// Update updates the game logic.
func (g *Game) Update() error {
	// Handle player movement
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.playerY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.playerY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += 2
	}
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Aztank Game!")

	// Draw the player tank using vector.DrawFilledRect
	playerColor := color.RGBA{255, 255, 255, 255} // White color
	vector.DrawFilledRect(screen, float32(g.playerX), float32(g.playerY), 20, 20, playerColor, false)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{playerX: screenWidth / 2, playerY: screenHeight / 2}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Aztank")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
