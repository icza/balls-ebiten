package engine

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/balls-ebiten/gfx"
)

// scene is used to present the world.
type scene struct {
	// e is the engine
	e *Engine
}

// newScene creates a new scene.
func newScene(e *Engine) *scene {
	s := &scene{
		e: e,
	}

	return s
}

// present paints the scene.
func (s *scene) present(screen *ebiten.Image) {
	ctx := gfx.NewContext(screen)

	// Paint background and frame:
	ctx.SetColorRGBA(100, 100, 100, 255)
	ctx.Rectangle(0, 0, s.e.w, s.e.h)

	// Paint balls:
	ctx.SetColorRGBA(200, 80, 0, 255)
	for _, b := range s.e.balls {
		s.paintBall(ctx, b)
	}

	s.paintOSD(ctx)

	s.paintGravity(ctx)
}

// paintOSD paints on-screen texts.
func (s *scene) paintOSD(ctx *gfx.Context) {
	if !s.e.osd {
		return
	}

	ctx.SetColorRGBA(200, 200, 100, 255)
	speed := 1.0
	if exp := s.e.speedExp; exp >= 0 {
		speed *= float64(int(1) << uint(exp))
	} else {
		speed /= float64(int(1) << uint(-exp))
	}

	phase := cmplx.Phase(s.e.gravity) + math.Pi*2 // Make sure it's positive
	degree := (720 - int(phase/math.Pi*180+0.5)) % 360

	items := []struct {
		keys   string
		format string
		param  interface{}
	}{
		{"R", "restart", nil},
		{"Q/X", "quit", nil},
		{"O", "OSD (on-screen display)", nil},
		{"S/s", "speed: %.2f", speed},
		{"A/a", "max # of balls: %2d", s.e.maxBalls},
		{"G/g", "abs gravity: %.2f", cmplx.Abs(s.e.gravity) / maxAbsGravity},
		{"T/t", "rotate gravity: %3d deg", degree},
		{"M/m", "min/max ball ratio: %.1f", float64(s.e.minMaxBallRatio) / 100},
	}

	col2x := func(col int) int { return col*210 + 10 }
	row2y := func(row int) int { return row*15 + 15 }

	// How many text columns fits on the screen?
	numCol := 0
	for col2x(numCol+1) < int(s.e.w) {
		numCol++
	}

	row, col := 0, 0
	for _, it := range items {
		params := []interface{}{"[" + it.keys + "]"}
		if it.param != nil {
			params = append(params, it.param)
		}
		text := fmt.Sprintf("%-5s "+it.format, params...)
		ctx.DrawString(text, col2x(col), row2y(row))

		col++
		if col >= numCol {
			row, col = row+1, 0
		}
	}
}

// paintBall paints the picture of a ball, a filled circle with 3D effects.
func (s *scene) paintBall(ctx *gfx.Context, b *ball) {
	// If performance becomes an issue, predraw on a texture,
	// cache it and just present the texture.
	x, y := int(real(b.pos)), int(imag(b.pos))

	// Fill circles going from outside
	gran := 10
	for i := 1; i <= gran; i++ {
		f := 1 - float64(i)/float64(gran+1)

		// Color is darker outside:
		col := func(c uint8) uint8 {
			return c - uint8(float64(c)*0.7*f)
		}

		ctx.SetColorRGBA(col(b.c.R), col(b.c.G), col(b.c.B), b.c.A)
		ctx.FillCircle(x, y, int(b.r*f))
	}

	ctx.SetColorRGBA(255, 255, 255, b.c.A)
	ctx.Point(x, y)
}

// paintGravity paints a gravity vector.
func (s *scene) paintGravity(ctx *gfx.Context) {
	const size = 70 // Pixel size of max gravity
	g := s.e.gravity * complex(float64(size)/maxAbsGravity, 0)

	x1, y1 := s.e.w-size-2, s.e.h-size-2
	x2, y2 := x1+int(real(g)), y1+int(imag(g))

	ctx.SetColorRGBA(50, 150, 255, 255)
	ctx.DrawLine(x1, y1, x2, y2)

	// Bottom of the arrow:
	v := g * 0.15i
	ctx.DrawLine(x1, y1, x1+int(real(v)), y1+int(imag(v)))
	v = g * -0.15i
	ctx.DrawLine(x1, y1, x1+int(real(v)), y1+int(imag(v)))

	// Head of the arrow:
	v = g * (-0.18 + 0.18i)
	ctx.DrawLine(x2, y2, x2+int(real(v)), y2+int(imag(v)))
	v = g * (-0.18 - 0.18i)
	ctx.DrawLine(x2, y2, x2+int(real(v)), y2+int(imag(v)))
}
