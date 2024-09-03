package player

import (
	"image/color"

	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Player represents the player's state.
type Player struct {
	X, Y float64
}

// NewPlayer creates a new player at the given position.
func NewPlayer(x, y float64) *Player {
	return &Player{X: x, Y: y}
}

// Update handles player movement.
func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += 2
	}
}

// Draw renders the player on the screen.
func (p *Player) Draw(screen *ebiten.Image) {
	playerColor := color.RGBA{255, 255, 255, 255} // White color
	vector.DrawFilledRect(screen, float32(p.X), float32(p.Y), utils.PlayerWidth, utils.PlayerHeight, playerColor, false)
}
