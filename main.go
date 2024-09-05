package main

import (
	"log"

	"github.com/FabioSM46/tank/internal/game"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
    game := game.NewGame()

    ebiten.SetWindowSize(utils.ScreenWidth, utils.ScreenHeight)
    ebiten.SetWindowTitle("GoTank Game")

    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
