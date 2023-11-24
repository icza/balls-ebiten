// Package main is the Bouncing balls (Ebiten) demo app.
package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	version  = "v1.0.0"
	name     = "Bouncing Balls Ebiten"
	homePage = "https://github.com/icza/balls-ebiten"
	title    = name + " " + version
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

const width, height = 800, 550

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	fmt.Println(title)
	fmt.Println("Home page:", homePage)

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)

	g := &Game{}

	if err := ebiten.RunGame(g); err != nil {
		log.Printf("RunGame() error: %v", err)
	}
}
