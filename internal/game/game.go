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
}

func NewGame() *Game {
    // Construct the paths to assets using the working directory
    player1ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank.png")
    player2ImagePath := filepath.Join(config.WorkingDir, "assets/images/tank.png")

    gameMap := NewGameMap()

    return &Game{
        player1: NewPlayer(100, 100, player1ImagePath, gameMap),  
        player2: NewAIPlayer(400, 300, player2ImagePath), 
        gameMap: gameMap,
    }
}

func (g *Game) Update() error {
    g.player1.Update()
    g.player2.UpdateAI(g.player1) // AI moves towards player1
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.gameMap.Draw(screen)
    g.player1.Draw(screen)
    g.player2.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return utils.ScreenWidth, utils.ScreenHeight
}
