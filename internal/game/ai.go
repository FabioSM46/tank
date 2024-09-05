package game

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// AIPlayer represents an AI-controlled player in the game.
type AIPlayer struct {
    X, Y   float64   // Position (exported fields)
    Speed  float64   // Movement speed
    Rotation float64      // Rotation angle in degrees
    Color  color.Color // AI color
    Image  *ebiten.Image // Tank image
}

// NewAIPlayer creates a new AI player instance.
func NewAIPlayer(x, y float64, imgPath string) *AIPlayer {
    img, _, err := ebitenutil.NewImageFromFile(imgPath) // Load the tank image
    if err != nil {
        log.Fatal(err)
    }

    return &AIPlayer{
        X:     x,
        Y:     y,
        Speed: 1.5, // Slower speed for AI
        Rotation: 0,
        Image: img,
    }
}

// UpdateAI moves the AI towards the player.
func (ai *AIPlayer) UpdateAI(player *Player) {
    // Calculate the direction towards the player
    dx := player.X - ai.X
    dy := player.Y - ai.Y
    angle := math.Atan2(dy, dx) * 180 / math.Pi // Convert radians to degrees

    // Update rotation to face the player
    ai.Rotation = angle

   // rad := ai.Rotation * math.Pi / 180 // Convert rotation angle to radians

    // Move the AI towards the player
    ai.X += 0 //ai.Speed * math.Cos(rad)
    ai.Y += 0//ai.Speed * math.Sin(rad)
}

// Draw renders the AI player on the screen.
func (ai *AIPlayer) Draw(screen *ebiten.Image) {
    opts := &ebiten.DrawImageOptions{}
    
    // Translate to the center of the image, rotate, then translate back
    imgWidth, imgHeight := ai.Image.Bounds().Dx(), ai.Image.Bounds().Dy()
    opts.GeoM.Translate(-float64(imgWidth)/2, -float64(imgHeight)/2)
    opts.GeoM.Rotate(ai.Rotation * math.Pi / 180) // Convert to radians
    opts.GeoM.Translate(float64(imgWidth)/2, float64(imgHeight)/2)
    
    // Translate to the position on the screen
    opts.GeoM.Translate(ai.X, ai.Y)
    screen.DrawImage(ai.Image, opts)
}