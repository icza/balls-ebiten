// Package gfx contains drawing utilities for a draw.Image destination.
package gfx

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Context provides primitive drawing operations on a draw.Image.
// All operations use the drawing color set by SetColor().
type Context struct {
	// dst is the image to draw on
	dst *ebiten.Image

	// col is the drawing color
	col color.Color
}

// NewContext returns a new Context initialized with black drawing color.
func NewContext(screen *ebiten.Image) *Context {
	return &Context{
		dst: screen,
		col: color.NRGBA{0, 0, 0, 255},
	}
}

// SetColor sets the drawing color.
func (ctx *Context) SetColor(col color.Color) {
	ctx.col = col
}

// SetColorRGBA sets the drawing color.
func (ctx *Context) SetColorRGBA(r, g, b, a byte) {
	ctx.SetColor(color.RGBA{r, g, b, a})
}

// Point draws a point.
func (ctx *Context) Point(x, y int) {
	ctx.dst.Set(x, y, ctx.col)
}

// Rectangle draws a rectangle.
func (ctx *Context) Rectangle(x1, y1, width, height int) {
	vector.StrokeRect(ctx.dst, float32(x1), float32(y1), float32(width), float32(height), 1, ctx.col, true)
}

// FillCircle draws a filled circle.
func (ctx *Context) FillCircle(x0, y0, rad int) {
	vector.DrawFilledCircle(ctx.dst, float32(x0), float32(y0), float32(rad), ctx.col, true)
}

// DrawString draws a string.
// The y coordinate is the bottom line of the text.
func (ctx *Context) DrawString(s string, x, y int) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  ctx.dst,
		Src:  image.NewUniform(ctx.col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(s)
}

// abs reutrns the absolute of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawLine draws a line.
func (ctx *Context) DrawLine(x0, y0, x1, y1 int) {
	vector.StrokeLine(ctx.dst, float32(x0), float32(y0), float32(x1), float32(y1), 1, ctx.col, true)
}
