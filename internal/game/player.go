package game

import (
	"image/color"
	"log"
	"math"

	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/toolchain/src/math/rand"
)

// Bullet represents a bullet shot by the player.
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

// Player represents a player in the game.
type Player struct {
	X, Y      float64     // Position
	Speed     float64     // Movement speed
	Rotation  float64     // Rotation angle in degrees
	Color     color.Color // Player color
	Image     *ebiten.Image // Tank image
	GameMap   *GameMap    // Reference to the game map for collision detection
	Health       int	// Health points
}

// NewPlayer creates a new player instance.
func NewPlayer(x, y float64, imgPath string, gameMap *GameMap) *Player {
	img, _, err := ebitenutil.NewImageFromFile(imgPath) // Load the tank image
	if err != nil {
		log.Fatal(err)
	}

	player := &Player{
		X:        x,
		Y:        y,
		Speed:    2.0,
		Rotation: 0, // Initial rotation angle
		Image:    img,
		GameMap:  gameMap, // Set the game map reference
		Health:   100,
	}

	// Ensure the player does not start on a wall tile
	playerWidth, playerHeight := img.Bounds().Dx(), img.Bounds().Dy()

	// Calculate the bounding box for collision checking
	left, top, right, bottom := player.GetBoundingBoxAt(x, y)

	// Check if the initial position is on a wall tile, and reposition if necessary
	for gameMap.CheckCollisionWithWall(left, top, right, bottom) {
		// Generate new random x, y positions within the map bounds
		player.X = float64(rand.Intn(MapWidth)*TileSize + playerWidth/2)
		player.Y = float64(rand.Intn(MapHeight)*TileSize + playerHeight/2)

		// Update the bounding box for the new position
		left, top, right, bottom = player.GetBoundingBoxAt(player.X, player.Y)
	}

	return player
}

// GetBoundingBox returns the bounding box of the player as (left, top, right, bottom)
func (p *Player) GetBoundingBox() (float64, float64, float64, float64) {
	imgWidth, imgHeight := p.Image.Bounds().Dx(), p.Image.Bounds().Dy()
	halfWidth := float64(imgWidth) / 2
	halfHeight := float64(imgHeight) / 2
	left := p.X - halfWidth
	top := p.Y - halfHeight
	right := p.X + halfWidth
	bottom := p.Y + halfHeight

	return left, top, right, bottom
}

func (p *Player) Update() {
	rad := p.Rotation * math.Pi / 180 // Convert rotation angle to radians

	// Calculate the potential new position
	newX, newY := p.X, p.Y
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

	// Calculate the player's bounding box at the new position
	newPx1, newPy1, newPx2, newPy2 := p.GetBoundingBoxAt(newX, newY)

	// Check for wall collision
	if !p.GameMap.CheckCollisionWithWall(newPx1, newPy1, newPx2, newPy2) {
		// Update the player position if no collision
		p.X = newX
		p.Y = newY
	}

	// Normalize the rotation angle to stay within 0-360 degrees
	if p.Rotation >= 360 {
		p.Rotation -= 360
	} else if p.Rotation < 0 {
		p.Rotation += 360
	}
}

// GetBoundingBoxAt returns the bounding box of the player at a specific position
func (p *Player) GetBoundingBoxAt(x, y float64) (float64, float64, float64, float64) {
	imgWidth, imgHeight := p.Image.Bounds().Dx(), p.Image.Bounds().Dy()
	halfWidth := float64(imgWidth) / 2
	halfHeight := float64(imgHeight) / 2

	offsetX, offsetY := 16.0, 16.0 // Adjust this value to move the bounding box
	margin := 5.0 // Adjust this value to make the bounding box narrower

	// Calculate the bounding box with the applied offsets
	left := x - halfWidth + margin + offsetX
	top := y - halfHeight + margin + offsetY
	right := left + float64(imgWidth) - 2*margin
	bottom := top + float64(imgHeight) - 2*margin

	return left, top, right, bottom
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

	// Draw the tank
	screen.DrawImage(p.Image, opts)

	// Draw the bounding box
	// Apply the same transformations to the bounding box
	bboxOpts := &ebiten.DrawImageOptions{}
	bboxOpts.GeoM.Translate(-float64(imgWidth)/2, -float64(imgHeight)/2)
	bboxOpts.GeoM.Rotate(p.Rotation * math.Pi / 180) // Convert to radians
	bboxOpts.GeoM.Translate(float64(imgWidth)/2, float64(imgHeight)/2)
	bboxOpts.GeoM.Translate(p.X, p.Y)

	left, top, right, bottom := p.GetBoundingBoxAt(p.X, p.Y)
	
	// Draw bounding box with red lines
	vector.StrokeLine(screen, float32(left), float32(top), float32(right), float32(top), 1, color.RGBA{255, 0, 0, 255}, false)   // Top line
	vector.StrokeLine(screen, float32(right), float32(top), float32(right), float32(bottom), 1, color.RGBA{255, 0, 0, 255}, false) // Right line
	vector.StrokeLine(screen, float32(right), float32(bottom), float32(left), float32(bottom), 1, color.RGBA{255, 0, 0, 255}, false) // Bottom line
	vector.StrokeLine(screen, float32(left), float32(bottom), float32(left), float32(top), 1, color.RGBA{255, 0, 0, 255}, false) // Left line
}

func (p *Player) TakeDamage(amount int) {
    p.Health -= amount
    if p.Health < 0 {
        p.Health = 0
    }
}