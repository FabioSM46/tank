package game

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
    GameMap  *GameMap // Reference to the game map for collision detection
}

// NewPlayer creates a new player instance.
func NewPlayer(x, y float64, imgPath string, gameMap *GameMap) *Player {
	img, _, err := ebitenutil.NewImageFromFile(imgPath) // Load the tank image
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		X:        x,
		Y:        y,
		Speed:    2.0,
		Rotation: 0, // Initial rotation angle
		Image:    img,
		GameMap:  gameMap, // Set the game map reference
	}
}


func (p *Player) Update() {
	rad := p.Rotation * math.Pi / 180 // Convert rotation angle to radians

	newX, newY := p.X, p.Y // Store the new position

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		newX += p.Speed * math.Cos(rad)
		newY += p.Speed * math.Sin(rad)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		newX -= p.Speed * math.Cos(rad)
		newY -= p.Speed * math.Sin(rad)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Rotation -= 2 // Rotate left
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Rotation += 2 // Rotate right
	}

    // Check for wall collision before moving
    if !p.GameMap.IsWall(newX, newY) {
        // Check if the new position is within the map boundaries
        if newX >= 0 && newX <= float64(MapWidth*TileSize) && newY >= 0 && newY <= float64(MapHeight*TileSize) {
            // Update the player position if no collision and within map bounds
            p.X = newX
            p.Y = newY
        }
    }

	// Normalize the rotation angle to stay within 0-360 degrees
	if p.Rotation >= 360 {
		p.Rotation -= 360
	} else if p.Rotation < 0 {
		p.Rotation += 360
	}
}

// Draw renders the player on the screen.
func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	// Translate to the center of the image
	imgWidth, imgHeight := p.Image.Bounds().Dx(), p.Image.Bounds().Dy()
	opts.GeoM.Translate(-float64(imgWidth)/2, -float64(imgHeight)/2)

	// Apply rotation
	opts.GeoM.Rotate(p.Rotation * math.Pi / 180) // Convert to radians

	// Translate back to the original position
	opts.GeoM.Translate(float64(imgWidth)/2, float64(imgHeight)/2)

	// Translate to the tank's position
	opts.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(p.Image, opts)
}