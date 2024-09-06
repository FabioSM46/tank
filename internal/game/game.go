package game

import (
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
				normalAngle := g.calculateWallNormal(bullet.X)
				bullet.Bounce(normalAngle)
			} else {
				// Remove the bullet after two bounces
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			}
		}

		// Remove bullets that are out of bounds
		if bullet.IsOutOfBounds() {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
		}
	}

	// Check for shooting input
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.ShootBullet(g.player1)
	}

	return nil
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
func (g *Game) calculateWallNormal(x float64) float64 {
    // Simple example: assuming vertical or horizontal walls
    // You can improve this function to calculate based on the actual wall angle
    if x <= 0 || x >= utils.ScreenWidth {
        return 90 // Normal for vertical walls
    }
    return 0 // Normal for horizontal walls
}