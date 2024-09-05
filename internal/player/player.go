package player

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Player represents a player in the game.
type Player struct {
    X, Y   float64   // Position (exported fields)
    Speed  float64   // Movement speed
    Rotation float64      // Rotation angle in degrees
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
        Rotation: 0, 
        Image: img,
    }
}

// Update handles player input and movement.
func (p *Player) Update() {
    rad := p.Rotation * math.Pi / 180 // Convert rotation angle to radians

    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        p.X += p.Speed * math.Cos(rad)
        p.Y += p.Speed * math.Sin(rad)
    }
    if ebiten.IsKeyPressed(ebiten.KeyDown) {
        p.X -= p.Speed * math.Cos(rad)
        p.Y -= p.Speed * math.Sin(rad)
    }
    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        p.Rotation -= 2 // Rotate left
    }
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        p.Rotation += 2 // Rotate right
    }
}

// Draw renders the player on the screen.
func (p *Player) Draw(screen *ebiten.Image) {
    opts := &ebiten.DrawImageOptions{}
    
    // Translate to the center of the image, rotate, then translate back
    imgWidth, imgHeight := p.Image.Bounds().Dx(), p.Image.Bounds().Dy()
    opts.GeoM.Translate(-float64(imgWidth)/2, -float64(imgHeight)/2)
    opts.GeoM.Rotate(p.Rotation * math.Pi / 180) // Convert to radians
    opts.GeoM.Translate(float64(imgWidth)/2, float64(imgHeight)/2)
    
    // Translate to the position on the screen
    opts.GeoM.Translate(p.X, p.Y)
    screen.DrawImage(p.Image, opts)
}