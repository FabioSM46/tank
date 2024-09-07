package game

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/toolchain/src/math/rand"
)

type AIPlayer struct {
    X, Y          float64     // Position
    Speed         float64     // Movement speed
    Rotation      float64     // Rotation angle in degrees
    Color         color.Color // AI color
    Image         *ebiten.Image // Tank image
    lastShotTime  float64     // Time of the last shot
    shootCooldown float64     // Time between shots (in seconds)
    GameMap       *GameMap    // Reference to the game map for collision detection
    stuckCounter  int         // Counter to detect if AI is stuck
    avoidAttempts int         // Number of attempts to avoid the obstacle
    Health int                // Health points
}

// NewAIPlayer creates a new AI player instance.
func NewAIPlayer(x, y float64, imgPath string, gameMap *GameMap) *AIPlayer {
    img, _, err := ebitenutil.NewImageFromFile(imgPath)
    if err != nil {
        log.Fatal(err)
    }

    aiPlayer := &AIPlayer{
        X:            x,
        Y:            y,
        Speed:        1,
        Rotation:     0,
        Image:        img,
        lastShotTime: 0,
        shootCooldown: 2.0, 
        GameMap:      gameMap,
        Health:       100,
    }

    // Ensure the AI player does not start on a wall tile
    imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()

    // Calculate the bounding box for collision checking
    left, top, right, bottom := aiPlayer.GetBoundingBoxAt(x, y)

    // Check if the initial position is on a wall tile, and reposition if necessary
    for gameMap.CheckCollisionWithWall(left, top, right, bottom) {
        // Generate new random x, y positions within the map bounds
        aiPlayer.X = float64(rand.Intn(MapWidth)*TileSize + imgWidth/2)
        aiPlayer.Y = float64(rand.Intn(MapHeight)*TileSize + imgHeight/2)

        // Update the bounding box for the new position
        left, top, right, bottom = aiPlayer.GetBoundingBoxAt(aiPlayer.X, aiPlayer.Y)
    }

    return aiPlayer
}

// GetBoundingBoxAt returns the bounding box of the AI player at a specific position.
func (ai *AIPlayer) GetBoundingBoxAt(x, y float64) (float64, float64, float64, float64) {
    imgWidth, imgHeight := ai.Image.Bounds().Dx(), ai.Image.Bounds().Dy()
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

func (ai *AIPlayer) UpdateAI(player *Player, currentTime float64, game *Game) {
    // Calculate the direction towards the player
    dx := player.X - ai.X
    dy := player.Y - ai.Y
    angle := math.Atan2(dy, dx) * 180 / math.Pi // Convert radians to degrees
    px1, py1, px2, py2 := player.GetBoundingBox()

    // Update rotation to face the player
    ai.Rotation = angle

    // Try moving towards the player
     ai.tryMove()

    // Check if the AI should shoot
    if currentTime-ai.lastShotTime >= ai.shootCooldown {
        if ai.GameMap.CheckLineOfSight(player.X, player.Y, ai.X, ai.Y) && CheckLineRectIntersection(ai.X, ai.Y, player.X, player.Y, px1, py1, px2, py2) {
            ai.ShootBullet(game)
            ai.lastShotTime = currentTime
        }
    }
}

// `tryMove` attempts to move the AI tank and adjusts its rotation to avoid obstacles.
func (ai *AIPlayer) tryMove() {
    rad := ai.Rotation * math.Pi / 180 // Convert rotation angle to radians

    // Calculate the potential new position
    newX, newY := ai.X+ai.Speed*math.Cos(rad), ai.Y+ai.Speed*math.Sin(rad)
    newAx1, newAy1, newAx2, newAy2 := ai.GetBoundingBoxAt(newX, newY)

    if !ai.GameMap.CheckCollisionWithWall(newAx1, newAy1, newAx2, newAy2) {
        // No collision; move the AI tank
        ai.X, ai.Y = newX, newY
        ai.resetStuckState()
        return
    }

    // Collision detected; try avoiding it
    ai.avoidObstacle()
}

func (ai *AIPlayer) avoidObstacle() {
    // Increment the stuck counter to track how many times avoidance has been attempted
    ai.stuckCounter++

    // Define the base avoidance angle
    baseAvoidanceAngle := 45.0

    // Dynamically increase the avoidance angle based on the stuckCounter
    avoidanceAngle := baseAvoidanceAngle + float64(ai.stuckCounter*15) // Increase by 15 degrees per attempt

    // Define a set of avoidance angles to try
    avoidanceAngles := []float64{avoidanceAngle, -avoidanceAngle, 2 * avoidanceAngle, -2 * avoidanceAngle}

    // Loop through each avoidance angle to check for a way around the obstacle
    for _, angle := range avoidanceAngles {
        // Rotate the tank by the current avoidance angle
        ai.Rotation += angle
        rad := ai.Rotation * math.Pi / 180

        // Calculate new potential position
        newX, newY := ai.X+ai.Speed*math.Cos(rad), ai.Y+ai.Speed*math.Sin(rad)
        newAx1, newAy1, newAx2, newAy2 := ai.GetBoundingBoxAt(newX, newY)

        // Check if this new position is free of collisions
        if !ai.GameMap.CheckCollisionWithWall(newAx1, newAy1, newAx2, newAy2) {
            // If the new position is free, update the position and reset counters
            ai.X, ai.Y = newX, newY
            ai.resetStuckState()
            return
        }
    }

    // If all attempts failed, increase avoidAttempts and change behavior
    ai.avoidAttempts++
    if ai.avoidAttempts > 3 { // After 3 failed attempts, try a random angle
        ai.Rotation += float64(45 + rand.Intn(90)) // Random angle between 45 and 135 degrees
        ai.avoidAttempts = 0                       // Reset the attempt counter
    }
}

// Reset the AI tank's stuck state counters
func (ai *AIPlayer) resetStuckState() {
    ai.stuckCounter = 0
    ai.avoidAttempts = 0
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

func (ai *AIPlayer) TakeDamage(amount int) {
    ai.Health -= amount
    if ai.Health < 0 {
        ai.Health = 0
    }
}