package game

import (
	"image/color"

	"github.com/FabioSM46/tank/internal/input"
	"github.com/FabioSM46/tank/internal/player"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    player1 *player.Player
    player2 *input.AIPlayer
    gameMap *GameMap
}

func NewGame() *Game {
    return &Game{
		player1: player.NewPlayer(100, 100, color.RGBA{0, 0, 255, 255}), // Blue player
        player2: input.NewAIPlayer(400, 300, color.RGBA{255, 0, 0, 255}), // Red AI player
        gameMap: NewGameMap(),
    }
}

func (g *Game) Update() error {
    g.player1.Update()
    g.player2.UpdateAI(g.player1) // Pass player1 to AI for behavior decisions
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
