// Package main is the Bouncing balls (Ebiten) demo app.
package main

import (
	"fmt"
	"log"
	"time"

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
	tl  *typeListener
	eng *engine.Engine // eng is the engine
}

func (g *Game) handleInputs() {
	eng := g.eng

	shift := ebiten.IsKeyPressed(ebiten.KeyShift)

	tl := g.tl
	tl.now = time.Now()

	if tl.Typed(ebiten.KeyS) {
		eng.ChangeSpeed(shift)
	}
	if tl.Typed(ebiten.KeyA) {
		eng.ChangeMaxBalls(shift)
	}
	if tl.Typed(ebiten.KeyM) {
		eng.ChangeMinMaxBallRatio(shift)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		eng.Restart()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		eng.ToggleOSD()
	}
	if tl.Typed(ebiten.KeyG) {
		eng.ChangeGravityAbs(shift)
	}
	if tl.Typed(ebiten.KeyT) {
		eng.RotateGravity(shift)
	}
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
		tl:  newTypeListener(),
		eng: engine.NewEngine(w, h),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Printf("RunGame() error: %v", err)
	}
}

type typeListener struct {
	now       time.Time
	lastTyped map[ebiten.Key]time.Time // Zero time is good for initial value
}

func newTypeListener() *typeListener {
	return &typeListener{
		lastTyped: map[ebiten.Key]time.Time{},
	}
}

func (tl *typeListener) Typed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(key) {
		if tl.now.Sub(tl.lastTyped[key]) > 200*time.Millisecond {
			tl.lastTyped[key] = tl.now
			return true
		}
	}
	return false
}
