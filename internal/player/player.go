package player

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Player represents a player in the game.
type Player struct {
    X, Y   float64   // Position (exported fields)
    Speed  float64   // Movement speed
    Color  color.Color // Player color
    Image  *ebiten.Image // Tank image
}

// NewPlayer creates a new player instance.
func NewPlayer(x, y float64, imgPath string) *Player {
    img, _, err := ebitenutil.NewImageFromFile(imgPath) // Load the tank image
    if err != nil {
        log.Fatal(err)
    }

    return &Player{
        X:     x,
        Y:     y,
        Speed: 2.0,
        Image: img,
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
    opts := &ebiten.DrawImageOptions{}
    opts.GeoM.Translate(p.X, p.Y)
    screen.DrawImage(p.Image, opts)
}
