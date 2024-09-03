package game

import (
	"github.com/FabioSM46/tank/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game represents the current game state.
type Game struct {
	Player *player.Player
}

// NewGame initializes and returns a new Game instance.
func NewGame(screenWidth, screenHeight int) *Game {
	return &Game{
		Player: player.NewPlayer(float64(screenWidth)/2, float64(screenHeight)/2),
	}
}

// Update updates the game logic.
func (g *Game) Update() error {
	g.Player.Update()
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Aztank Game!")
	g.Player.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
