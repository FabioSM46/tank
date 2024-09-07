package game

import (
	"image/color"
	"math"

	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bullet struct {
	X, Y        float64 // Position
	VelX, VelY  float64 // Velocity components
	Rotation    float64 // Direction in degrees
	BounceCount int     // Number of times the bullet has bounced
}

// NewBullet creates a new bullet instance starting slightly in front of the tank.
func NewBullet(player *Player) *Bullet {
	// Calculate the center of the tank
	imgWidth, imgHeight := player.Image.Bounds().Dx(), player.Image.Bounds().Dy()
	centerX := player.X + float64(imgWidth)/2
	centerY := player.Y + float64(imgHeight)/2

	// Calculate the starting position in front of the tank
	rad := player.Rotation * math.Pi / 180 // Convert rotation angle to radians
	offsetX := 20 * math.Cos(rad)          // Distance in front of the tank
	offsetY := 20 * math.Sin(rad)          // Distance in front of the tank

	// Set initial speed and velocity components
	speed := 5.0
	velX := speed * math.Cos(rad)
	velY := speed * math.Sin(rad)

	return &Bullet{
		X:          centerX + offsetX,
		Y:          centerY + offsetY,
		VelX:       velX,
		VelY:       velY,
		Rotation:   player.Rotation,
		BounceCount: 0,
	}
}

// NewBulletFromAI creates a bullet instance for the AI
func NewBulletFromAI(ai *AIPlayer) *Bullet {
    // Calculate the center of the AI tank
    imgWidth, imgHeight := ai.Image.Bounds().Dx(), ai.Image.Bounds().Dy()
    centerX := ai.X + float64(imgWidth)/2
    centerY := ai.Y + float64(imgHeight)/2

    // Calculate the starting position in front of the AI tank
    rad := ai.Rotation * math.Pi / 180 // Convert rotation angle to radians
    offsetX := 20 * math.Cos(rad)      // Distance in front of the tank
    offsetY := 20 * math.Sin(rad)      // Distance in front of the tank

    // Set initial speed and velocity components
    speed := 5.0
    velX := speed * math.Cos(rad)
    velY := speed * math.Sin(rad)

    return &Bullet{
        X:          centerX + offsetX,
        Y:          centerY + offsetY,
        VelX:       velX,
        VelY:       velY,
        Rotation:   ai.Rotation,
        BounceCount: 0,
    }
}

// Update moves the bullet according to its velocity.
func (b *Bullet) Update() {
	b.X += b.VelX
	b.Y += b.VelY
}

// Bounce changes the bullet's direction when it hits a wall.
func (b *Bullet) Bounce(normalAngle float64) {
	// Convert the normal angle to radians
	normalRad := normalAngle * math.Pi / 180

	// Reflect the bullet's velocity vector over the normal
	dot := b.VelX*math.Cos(normalRad) + b.VelY*math.Sin(normalRad)
	b.VelX -= 2 * dot * math.Cos(normalRad)
	b.VelY -= 2 * dot * math.Sin(normalRad)

	b.BounceCount++
}

// IsOutOfBounds checks if the bullet is outside the game area.
func (b *Bullet) IsOutOfBounds() bool {
	return b.X < 0 || b.Y < 0 || b.X > utils.ScreenWidth || b.Y > utils.ScreenHeight
}

// Draw renders the bullet on the screen.
func (b *Bullet) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(b.X), float32(b.Y), 2, color.RGBA{255, 0, 0, 255}, false)
}