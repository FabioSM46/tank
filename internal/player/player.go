package player

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Player represents a player in the game.
type Player struct {
    X, Y   float64   // Position (exported fields)
    Speed  float64   // Movement speed
    Color  color.Color // Player color
}

// NewPlayer creates a new player instance.
func NewPlayer(x, y float64, color color.Color) *Player {
    return &Player{
        X:     x,
        Y:     y,
        Speed: 2.0,
        Color: color,
    }
}

// Update handles player input and movement.
func (p *Player) Update() {
    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        p.Y -= p.Speed
    }
    if ebiten.IsKeyPressed(ebiten.KeyDown) {
        p.Y += p.Speed
    }
    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        p.X -= p.Speed
    }
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        p.X += p.Speed
    }
}

// Draw renders the player on the screen.
func (p *Player) Draw(screen *ebiten.Image) {
    vector.DrawFilledRect(screen, float32(p.X), float32(p.Y), 20, 20, p.Color, false)
}
