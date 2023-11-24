// Package main is the Bouncing balls (Ebiten) demo app.
package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/balls-ebiten/engine"
)

const (
	version  = "v1.0.0"
	name     = "Bouncing Balls Ebiten"
	homePage = "https://github.com/icza/balls-ebiten"
	title    = name + " " + version
)

type Game struct {
	eng *engine.Engine // eng is the engine
}

func (g *Game) Update() error {
	g.eng.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.eng.Present(screen)
}

const w, h = 800, 550

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w, h
}

func main() {
	fmt.Println(title)
	fmt.Println("Home page:", homePage)

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(w, h)

	g := &Game{
		eng: engine.NewEngine(w, h),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Printf("RunGame() error: %v", err)
	}
}
