package main

import (
	"log"

	"github.com/FabioSM46/tank/internal/game"
	"github.com/FabioSM46/tank/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)



func main() {
	g := game.NewGame(utils.ScreenWidth, utils.ScreenHeight)
	ebiten.SetWindowSize(utils.ScreenWidth, utils.ScreenHeight)
	ebiten.SetWindowTitle("Aztank")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
