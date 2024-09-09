package main

import (
	"log"

	"github.com/FabioSM46/tank/internal/game"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	settings := utils.LoadWindowSettings()
	ebiten.SetWindowPosition(settings.X, settings.Y)
	ebiten.SetWindowSize(utils.ScreenWidth, utils.ScreenHeight)
	ebiten.SetWindowTitle("GoTank Game")
}

func main() {
	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	defer func() {
		x, y := ebiten.WindowPosition()
		utils.SaveWindowSettings(utils.WindowSettings{X: x, Y: y})
	}()
}
