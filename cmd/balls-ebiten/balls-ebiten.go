// Package main is the Bouncing balls (Ebiten) demo app.
package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (g *Game) handleInputs() {
	eng := g.eng

	shift := ebiten.IsKeyPressed(ebiten.KeyShift)

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		eng.ChangeSpeed(shift)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		eng.ChangeMaxBalls(shift)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		eng.ChangeMinMaxBallRatio(shift)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		eng.Restart()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		eng.ToggleOSD()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		eng.ChangeGravityAbs(shift)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		eng.RotateGravity(shift)
	}

	// TODO add to screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	g.handleInputs()

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
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{
		eng: engine.NewEngine(w, h),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Printf("RunGame() error: %v", err)
	}
}
