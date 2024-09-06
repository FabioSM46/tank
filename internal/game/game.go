package game

import (
	"math"
	"path/filepath"

	"github.com/FabioSM46/tank/config"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    player1 *Player
    player2 *AIPlayer
    gameMap *GameMap
    bullets []*Bullet 
}

func NewGame() *Game {
    // Construct the paths to assets using the working directory
    player1ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank.png")
    player2ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank.png")

    gameMap := NewGameMap()

    return &Game{
        player1: NewPlayer(100, 100, player1ImagePath, gameMap),  
        player2: NewAIPlayer(400, 300, player2ImagePath, gameMap),
        gameMap: gameMap,
        bullets: []*Bullet{},
    }
}

func (g *Game) Update() error {
	g.player1.Update()
	g.player2.UpdateAI(g.player1)

	// Update bullets
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		bullet.Update()
	
		// Check for wall collision
		if g.gameMap.CheckCollisionWithWall(bullet.X, bullet.Y, bullet.X, bullet.Y) {
			if bullet.BounceCount < 2 {
				// Calculate normal angle of the wall for bouncing
				normalAngle := g.calculateWallNormal(bullet.X, bullet.Y)
				bullet.Bounce(normalAngle)
			} else {
				// Remove the bullet after two bounces
				g.bullets = removeBullet(g.bullets, i)
				continue // Continue to the next iteration to avoid processing an invalid index
			}
		}
	
		// Remove bullets that are out of bounds
		if bullet.IsOutOfBounds() {
			g.bullets = removeBullet(g.bullets, i)
		}
	}
	

	// Check for shooting input
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.ShootBullet(g.player1)
	}

	return nil
}

func removeBullet(bullets []*Bullet, i int) []*Bullet {
    return append(bullets[:i], bullets[i+1:]...)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gameMap.Draw(screen)
	g.player1.Draw(screen)
	g.player2.Draw(screen)

	// Draw bullets
	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}
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
