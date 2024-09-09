package game

import (
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/FabioSM46/tank/config"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player1      *Player
	player2      *AIPlayer
	gameMap      *GameMap
	bullets      []*Bullet
	isPaused     bool
	winner       string
	player1Score int
	player2Score int
}

func NewGame() *Game {
	// Construct the paths to assets using the working directory
	player1ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank.png")
	player2ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank-ai.png")

	gameMap := NewGameMap()

	return &Game{
		player1:      NewPlayer(100, 100, player1ImagePath, gameMap, 100),
		player2:      NewAIPlayer(400, 300, player2ImagePath, gameMap, 100),
		gameMap:      gameMap,
		bullets:      []*Bullet{},
		isPaused:     false,
		winner:       "",
		player1Score: 0,
		player2Score: 0,
	}
}

func (g *Game) Update() error {
	if g.isPaused {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.RestartGame()
		}
		return nil
	}

	currentTime := float64(time.Now().UnixNano()) / float64(time.Second)

	g.player1.Update(currentTime, g)
	g.player2.UpdateAI(g.player1, currentTime, g)

	// Update bullets
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		bullet.Update()

		// Check for bullet collision with walls
		if g.gameMap.CheckCollisionWithWall(bullet.X, bullet.Y, bullet.X, bullet.Y) {
			if bullet.BounceCount < 2 {
				normalAngle := g.calculateWallNormal(bullet.X, bullet.Y)
				bullet.Bounce(normalAngle)
			} else {
				g.bullets = removeBullet(g.bullets, i)
				continue
			}
		}

		// Check for bullet collision with tanks, ignoring bullets from the same owner
		if bullet.Owner != g.player1 && g.checkBulletCollisionWithTank(bullet, g.player1) {
			g.player1.TakeDamage(10, g) // Pass the game instance to handle death
			g.bullets = removeBullet(g.bullets, i)
			continue
		}

		if bullet.Owner != g.player2 && g.checkBulletCollisionWithTank(bullet, g.player2) {
			g.player2.TakeDamage(10, g) // Pass the game instance to handle death
			g.bullets = removeBullet(g.bullets, i)
			continue
		}

		// Remove bullets that are out of bounds
		if bullet.IsOutOfBounds() {
			g.bullets = removeBullet(g.bullets, i)
		}
	}

	return nil
}

func (g *Game) checkBulletCollisionWithTank(bullet *Bullet, tank interface{}) bool {
	bulletSize := 5.0 // Define bullet size for collision detection
	var tankX, tankY, tankWidth, tankHeight float64

	// Type assertion to determine if it's a Player or AIPlayer
	switch t := tank.(type) {
	case *Player:
		tankX, tankY = t.X, t.Y
		tankWidth, tankHeight = 32.0, 32.0 // Size for player tank
	case *AIPlayer:
		tankX, tankY = t.X, t.Y
		tankWidth, tankHeight = 32.0, 32.0 // Size for AI tank
	default:
		return false // If not a known type, no collision
	}

	// Check if the bullet is within the bounds of the tank
	return bullet.X > tankX-bulletSize && bullet.X < tankX+tankWidth+bulletSize &&
		bullet.Y > tankY-bulletSize && bullet.Y < tankY+tankHeight+bulletSize
}

func removeBullet(bullets []*Bullet, i int) []*Bullet {
	return append(bullets[:i], bullets[i+1:]...)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gameMap.Draw(screen)
	g.player1.Draw(screen)
	g.player2.Draw(screen)

	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}

	g.drawHealthDisplay(screen)

	if g.isPaused {
		// Screen dimensions
		screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

		// Create a temporary image to measure text size
		message := "Paused"
		scoreText := "Player Score: " + strconv.Itoa(g.player1Score) + " | AI Score: " + strconv.Itoa(g.player2Score)
		restartText := "Press Enter to Restart"

		// Draw the "Paused" message
		messageWidth := len(message) * 8 // Character width
		messageHeight := 16              // Character height
		messageX := (screenWidth - messageWidth) / 2
		messageY := (screenHeight - messageHeight*3) / 2

		ebitenutil.DebugPrintAt(screen, message, messageX, messageY)

		// Draw the score text
		scoreX := (screenWidth - len(scoreText)*8) / 2
		scoreY := messageY + messageHeight

		ebitenutil.DebugPrintAt(screen, scoreText, scoreX, scoreY)

		// Draw the restart instruction
		restartX := (screenWidth - len(restartText)*8) / 2
		restartY := scoreY + messageHeight

		ebitenutil.DebugPrintAt(screen, restartText, restartX, restartY)
	}
}

func (g *Game) drawHealthDisplay(screen *ebiten.Image) {
	// Define text positions within the black band
	bandHeight := 30
	bandY := utils.ScreenHeight - bandHeight
	bandX := utils.ScreenWidth

	// Draw player 1 health
	healthText := "Player Health: " + strconv.Itoa(g.player1.Health)
	ebitenutil.DebugPrintAt(screen, healthText, 10, bandY+10)

	// Draw player 2 health
	healthText = "AI Health: " + strconv.Itoa(g.player2.Health)
	ebitenutil.DebugPrintAt(screen, healthText, bandX-100, bandY+10)

	// Draw scores
	scoreText := "Player Score: " + strconv.Itoa(g.player1Score) + " | AI Score: " + strconv.Itoa(g.player2Score)
	ebitenutil.DebugPrintAt(screen, scoreText, (bandX/2)-100, bandY+10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return utils.ScreenWidth, utils.ScreenHeight
}

// ShootBullet creates a new bullet from the player's current position and rotation
func (g *Game) ShootBullet(p *Player) {
	bullet := NewBullet(p)
	g.bullets = append(g.bullets, bullet)
}

// calculateWallNormal returns the normal angle of the wall at the collision point.
func (g *Game) calculateWallNormal(bulletX, bulletY float64) float64 {
	tileSize := float64(TileSize)

	// Determine the tile the bullet is in
	tileX := int(bulletX / tileSize)
	tileY := int(bulletY / tileSize)

	// Check if we are out of bounds
	if tileX < 0 || tileX >= MapWidth || tileY < 0 || tileY >= MapHeight {
		return 0 // Assume horizontal wall normal for out-of-bounds
	}

	// Calculate the center of the tile
	tileCenterX := float64(tileX)*tileSize + tileSize/2
	tileCenterY := float64(tileY)*tileSize + tileSize/2

	// Calculate the differences between the bullet and the tile center
	dx := bulletX - tileCenterX
	dy := bulletY - tileCenterY

	// Determine the primary axis of collision
	if math.Abs(dx) > math.Abs(dy) {
		// Collision is more horizontal
		if dx > 0 {
			return 180 // Right side of the tile
		}
		return 0 // Left side of the tile
	} else {
		// Collision is more vertical
		if dy > 0 {
			return 90 // Bottom side of the tile
		}
		return 270 // Top side of the tile
	}
}

func (g *Game) HandlePlayerDeath(player *Player) {
	// Check if the player is the winner
	g.winner = "Player 1"
	g.player2Score++

	// Pause the game
	g.isPaused = true
}

func (g *Game) HandleAIDeath(ai *AIPlayer) {
	g.winner = "AI"
	g.player1Score++

	g.isPaused = true
}

func (g *Game) RestartGame() {
	// Reset game state
	g.player1 = NewPlayer(100, 100, filepath.Join(config.WorkingDir, "assets/images/tank.png"), g.gameMap, 100)
	g.player2 = NewAIPlayer(400, 300, filepath.Join(config.WorkingDir, "assets/images/tank.png"), g.gameMap, 100)
	g.bullets = []*Bullet{}
	g.isPaused = false
	g.winner = ""
}
