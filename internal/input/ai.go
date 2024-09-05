package input

import (
	"image/color"
	"math"

	"github.com/FabioSM46/tank/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// AIPlayer represents an AI-controlled player in the game.
type AIPlayer struct {
    X, Y   float64   // Position (exported fields)
    Speed  float64   // Movement speed
    Color  color.Color // AI color
}

// NewAIPlayer creates a new AI player instance.
func NewAIPlayer(x, y float64, color color.Color) *AIPlayer {
    return &AIPlayer{
        X:     x,
        Y:     y,
        Speed: 1.5, // Slower speed for AI
        Color: color,
    }
}

// UpdateAI moves the AI towards the player.
func (ai *AIPlayer) UpdateAI(player *player.Player) {
    // Calculate direction towards the player.
    dx := player.X - ai.X
    dy := player.Y - ai.Y
    distance := math.Sqrt(dx*dx + dy*dy)

    // Move the AI towards the player.
    if distance > 0 {
        ai.X += ai.Speed * dx / distance
        ai.Y += ai.Speed * dy / distance
    }
}

// Draw renders the AI player on the screen.
func (ai *AIPlayer) Draw(screen *ebiten.Image) {
    vector.DrawFilledRect(screen, float32(ai.X), float32(ai.Y), 20, 20, ai.Color, false)
}
