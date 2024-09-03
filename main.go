package main

import (
	"log"

	"github.com/FabioSM46/tank/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func main() {
	g := game.NewGame(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Aztank")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
