package game

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type AIPlayer struct {
    X, Y          float64   // Position
    Speed         float64   // Movement speed
    Rotation      float64   // Rotation angle in degrees
    Color         color.Color // AI color
    Image         *ebiten.Image // Tank image
    lastShotTime  float64   // Time of the last shot
    shootCooldown float64   // Time between shots (in seconds)
}

func NewAIPlayer(x, y float64, imgPath string, gameMap *GameMap) *AIPlayer {
    img, _, err := ebitenutil.NewImageFromFile(imgPath)
    if err != nil {
        log.Fatal(err)
    }

    return &AIPlayer{
        X:            x,
        Y:            y,
        Speed:        1.5,
        Rotation:     0,
        Image:        img,
        lastShotTime: 0,
        shootCooldown: 2.0, // 2 seconds cooldown
    }
}

// UpdateAI moves the AI towards the player and makes it shoot when a valid hit is possible.
func (ai *AIPlayer) UpdateAI(player *Player, currentTime float64, game *Game) {
    // Calculate the direction towards the player
    dx := player.X - ai.X
    dy := player.Y - ai.Y
    angle := math.Atan2(dy, dx) * 180 / math.Pi // Convert radians to degrees

    // Update rotation to face the player
    ai.Rotation = angle

    rad := ai.Rotation * math.Pi / 180 // Convert rotation angle to radians

    // Move the AI towards the player
    ai.X += ai.Speed * math.Cos(rad)
    ai.Y += ai.Speed * math.Sin(rad)

    // Check if the AI should shoot
    if currentTime-ai.lastShotTime >= ai.shootCooldown {
        if ai.CanHitPlayer(player) {
            ai.ShootBullet(game)
            ai.lastShotTime = currentTime
        }
    }
}

// CanHitPlayer checks if the AI can hit the player
func (ai *AIPlayer) CanHitPlayer(player *Player) bool {
    // Get the bounding box of the player
    px1, py1, px2, py2 := player.GetBoundingBox()
    
    // Use the line of sight check
    return CheckLineRectIntersection(ai.X, ai.Y, ai.X+math.Cos(ai.Rotation*math.Pi/180)*1000, ai.Y+math.Sin(ai.Rotation*math.Pi/180)*1000, px1, py1, px2, py2)
}


// ShootBullet creates a new bullet from the AI's current position and rotation
func (ai *AIPlayer) ShootBullet(game *Game) {
    bullet := NewBulletFromAI(ai)
    game.bullets = append(game.bullets, bullet)
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


// NewBulletForAI creates a new bullet instance starting slightly in front of the AI tank.
func NewBulletForAI(ai *AIPlayer) *Bullet {
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

func CheckLineRectIntersection(x1, y1, x2, y2, left, top, right, bottom float64) bool {
    // Check if the line segment intersects with any of the rectangle's edges
    if LineIntersect(x1, y1, x2, y2, left, top, left, bottom) ||
        LineIntersect(x1, y1, x2, y2, left, top, right, top) ||
        LineIntersect(x1, y1, x2, y2, right, top, right, bottom) ||
        LineIntersect(x1, y1, x2, y2, left, bottom, right, bottom) {
        return true
    }
    return false
}

// LineIntersect checks if two line segments intersect
func LineIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
    denom := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
    if denom == 0 {
        return false // Lines are parallel
    }

    ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / denom
    ub := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / denom

    return ua >= 0 && ua <= 1 && ub >= 0 && ub <= 1
}
