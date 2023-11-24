package engine

import (
	"math"
	"math/cmplx"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// physicsPeriod is the period of model update
	physicsPeriod = time.Millisecond * 2 // 500 / sec

	// ballSpawnPeriod is the ball spawning period
	ballSpawnPeriod = time.Second

	// near defines how near we let balls close boundaries and each other (3=touch)
	near = 3

	// minSpeedExp is the min allowed speed exponent value for the simulation speed
	minSpeedExp = -5

	// maxSpeedExp is the max allowed speed exponent value for the simulation speed
	maxSpeedExp = 3

	// maxAbsGravity is the max absolute value of the gravity
	maxAbsGravity = 2000
)

// Engine is the simulation engine.
// Contains the model, controls the simulation and presents it on the screen
// (via the scene).
type Engine struct {
	// w and h are the width and height of the world
	w, h int

	// lastUpdate is the last update timestamp
	lastUpdate time.Time

	// lastSpawned is the last ball spawning timestamp
	lastSpawned time.Time

	// balls of the simulation
	balls []*ball

	// scene is used to present the world
	scene *scene

	// speedExp is the (relative) speed exponent of the simulation: 2^speedExp
	// 0 being the normal (1x), 1 being 2x, 2 being 4x, -1 being 1/2 etc.
	speedExp int

	// maxBalls is the max number of balls
	maxBalls int

	// minMaxBallRatio is the ratio of the possible min and max ball radius (%)
	minMaxBallRatio int

	// osd tells if OSD is visible
	osd bool

	// gravity is the gravity vector
	gravity complex128
}

// NewEngine creates a new Engine.
func NewEngine(w, h int) *Engine {
	e := &Engine{
		w:               w,
		h:               h,
		lastUpdate:      time.Now(),
		maxBalls:        15,
		minMaxBallRatio: 60,
		osd:             true,
		gravity:         0 + 600i,
	}
	e.scene = newScene(e)

	return e
}

// update updates (recalculates) the world.
// It does it incrementally until engine state reaches the current timestamp.
func (e *Engine) Present(screen *ebiten.Image) {
	e.scene.present(screen)
}

// update updates (recalculates) the world.
// It does it incrementally until engine state reaches the current timestamp.
func (e *Engine) Update() {
	now := time.Now()
	dtMax := now.Sub(e.lastUpdate)

	// Protection against "timeouts":
	// If excessive time elapsed since last call, skip the "missing" time
	// (typical causes include computer sleep and extreme system load).
	// presentPeriod is the period of the scene presentation
	const updateTimeout = 400 * time.Millisecond
	if dtMax > updateTimeout {
		dtMax = updateTimeout
	}

	// dt might be "big", much bigger than physics period, in which case
	// the balls might move huge distances. To avoid any "unexpected" states,
	// do granular movement.

	if len(e.balls) > e.maxBalls {
		e.balls = e.balls[:e.maxBalls]
	}
	if len(e.balls) < e.maxBalls && now.Sub(e.lastSpawned) > ballSpawnPeriod {
		e.spawnBall()
		e.lastSpawned = now
	}

	if exp := e.speedExp; exp >= 0 {
		dtMax *= 1 << uint(exp)
	} else {
		dtMax /= 1 << uint(-exp)
	}

	// We always pass physicsPeriod to recalcInternal(), which means
	// we get the exact same result no matter the speed.
	for dt := time.Duration(0); dt < dtMax; dt += physicsPeriod {
		e.updateUnit(physicsPeriod)
	}

	e.lastUpdate = now
}

// updateUnit recalculates the world for a time unit (physicsPeriod),
// in one step.
func (e *Engine) updateUnit(dt time.Duration) {
	dtSec := dt.Seconds()

	for i, b := range e.balls {
		oldX, oldY := real(b.pos), imag(b.pos)
		b.update(dtSec)
		x, y := real(b.pos), imag(b.pos)

		collision := false

		// Check if world boundaries are reached, and bounce back if so:
		if x < b.r-near || x >= float64(e.w)-b.r+near {
			b.v = complex(-real(b.v), imag(b.v))
			collision = true
		}
		if y < b.r-near || y >= float64(e.h)-b.r+near {
			b.v = cmplx.Conj(b.v)
			collision = true
		}

		// Check collision with other balls:
		x1, y1, x2, y2 := x-b.r, y-b.r, x+b.r, y+b.r
		for j, b2 := range e.balls {
			if i == j {
				continue
			}

			// Fast check: enclosing rectangle collision:
			b2x, b2y := real(b2.pos), imag(b2.pos)
			if x2 < b2x-b2.r ||
				x1 > b2x+b2.r ||
				y2 < b2y-b2.r ||
				y1 > b2y+b2.r {
				continue // enclosing rectangles do not intersect
			}

			// Exact check:
			if b1 := b; cmplx.Abs(b1.pos-b2.pos) < b1.r+b2.r-2*near {
				collision = true
				// Algo description: https://en.wikipedia.org/wiki/Elastic_collision
				// New velocities:
				dpos := b1.pos - b2.pos
				common := 2 / (b1.m + b2.m) / abssq(dpos)

				v1 := b1.v - common*b2.m*sprod(b1.v-b2.v, +dpos)*+dpos
				v2 := b2.v - common*b1.m*sprod(b2.v-b1.v, -dpos)*-dpos

				b1.v, b2.v = v1, v2
			}
		}

		if collision {
			b.pos = complex(oldX, oldY)
		}
	}
}

// scalar production:
func sprod(a, b complex128) complex128 {
	return complex(real(a)*real(b)+imag(a)*imag(b), 0)
}

// abs then square:
func abssq(a complex128) complex128 {
	x := cmplx.Abs(a)
	return complex(x*x, 0)
}

// spawnBall spawns a new ball.
func (e *Engine) spawnBall() {
	for i := 0; i < 100; i++ { // Retry 100 times if needed
		b := newBall(e)

		// Check collision with other balls:
		x, y, R := real(b.pos), imag(b.pos), 2*b.r // 2*r: leave bigger space than needed
		x1, y1, x2, y2 := x-R, y-R, x+R, y+R

		collides := false
		for _, b2 := range e.balls {
			// Fast check is enough for us: enclosing rectangle collision:
			b2x, b2y := real(b2.pos), imag(b2.pos)
			if x2 < b2x-b2.r ||
				x1 > b2x+b2.r ||
				y2 < b2y-b2.r ||
				y1 > b2y+b2.r {
				continue // enclosing rectangles do not intersect
			}
			collides = true
			break
		}

		if !collides {
			e.balls = append(e.balls, b)
			break
		}
	}
}

// Restart restarts the simulation: removes all balls.
func (e *Engine) Restart() {
	e.balls = nil
}

// ChangeSpeed changes the speed of the simulation.
// Doubles it if up is true, else halves it.
func (e *Engine) ChangeSpeed(up bool) {
	if up && e.speedExp < maxSpeedExp {
		e.speedExp++
	}
	if !up && e.speedExp > minSpeedExp {
		e.speedExp--
	}
}

// ChangeMinMaxBallRatio changes the min-max ball ratio.
// Adds +/- 0.1.
func (e *Engine) ChangeMinMaxBallRatio(up bool) {
	if up && e.minMaxBallRatio < 100 {
		e.minMaxBallRatio += 10
	}
	if !up && e.minMaxBallRatio > 10 {
		e.minMaxBallRatio -= 10
	}
}

// ChangeMaxBalls changes the max number of balls.
// Adds +/- 1.
func (e *Engine) ChangeMaxBalls(up bool) {
	if up && e.maxBalls < 50 {
		e.maxBalls++
	}
	if !up && e.maxBalls > 1 {
		e.maxBalls--
	}
}

// ToggleOSD toggles the OSD.
func (e *Engine) ToggleOSD() {
	e.osd = !e.osd
}

// ChangeGravityAbs changes the absolute value of the gravity.
// Adds +/- 100.
func (e *Engine) ChangeGravityAbs(up bool) {
	// convert to int so we always see "nice" results
	oldR := int(cmplx.Abs(e.gravity))
	r := oldR
	if up {
		if r += 100; r > maxAbsGravity {
			r = maxAbsGravity
		}
	} else {
		if r -= 100; r < 2 { // Don't let below 1 else we lose direction info
			r = 1
		}
	}
	e.gravity *= complex(float64(r)/float64(oldR), 0)
}

// RotateGravity rotates the gravity vector.
// Rotates +/- 10 degrees.
func (e *Engine) RotateGravity(counterClockwise bool) {
	// convert to int so we always see "nice" results
	r, theta := cmplx.Polar(e.gravity)
	degree := int((theta+math.Pi*2)/math.Pi*180 + 0.5)
	if counterClockwise {
		degree += 10
	} else {
		degree -= 10
	}
	e.gravity = cmplx.Rect(r, float64(degree)/180*math.Pi)
}
