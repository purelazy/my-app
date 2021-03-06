package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// ------------------ GLOBALS/STRCUCTS ----------

type col struct {
	r, g, b byte
}

// Window width and height
const w, h int32 = 800, 600

// Pixels
var pixels = make([]byte, w*h*4)

// ------------------ UTILS ---------------------

func ifNilPanic(err error) {
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func fillScreen(c col) {
	// col{byte(x % mod), byte(y % mod), 0}
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			putPixel(x, y, c)
		}
	}
}

// Abs returns the absolute value of x.
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func absInt(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func putPixel(x, y int32, c col) {
	i := (y*w + x) * 4

	if i < int32(len(pixels))-4 && i >= 0 {
		pixels[i] = c.r
		pixels[i+1] = c.g
		pixels[i+2] = c.b
	}
}

func bresenham(x0, y0, x1, y1 int32, c col) {
	if absInt(y1-y0) < absInt(x1-x0) {
		if x0 > x1 {
			bresenhamLo(x1, y1, x0, y0, c)
		} else {
			bresenhamLo(x0, y0, x1, y1, c)
		}
	} else {
		if y0 > y1 {
			bresenhamHi(x1, y1, x0, y0, c)
		} else {
			bresenhamHi(x0, y0, x1, y1, c)
		}
	}
}

func bresenhamLo(x0, y0, x1, y1 int32, c col) {
	dx := x1 - x0
	dy := y1 - y0
	var yi int32 = 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}

	D := 2*dy - dx
	y := y0

	for x := x0; x <= x1; x++ {

		putPixel(x, y, c)

		if D > 0 {
			y += yi
			D -= 2 * dx
		}
		D += 2 * dy
	}
}

func bresenhamHi(x0, y0, x1, y1 int32, c col) {
	dx := x1 - x0
	dy := y1 - y0
	var xi int32 = 1

	if dx < 0 {
		xi = -1
		dx = -dx
	}

	D := 2*dx - dy
	x := x0

	for y := y0; y <= y1; y++ {

		putPixel(x, y, c)

		if D > 0 {
			x += xi
			D -= 2 * dy
		}
		D += 2 * dx
	}
}

func dda(x1i, y1i, x2i, y2i int32, c col) {
	x1, y1, x2, y2 := float32(x1i), float32(y1i), float32(x2i), float32(y2i)
	dx := x2 - x1
	dy := y2 - y1

	steps := max(abs(dx), abs(dy))

	xinc := dx / steps
	yinc := dy / steps

	for i := 1; i <= int(steps); i++ {
		putPixel(int32(x1), int32(y1), c)
		x1 += xinc
		y1 += yinc
	}
}

// ----------------------------------------------------------------------------

type vec2 struct {
	x, y float32
}

type ball struct {
	pos vec2
	vel vec2
	rad int32
	col col
}

func (ball *ball) draw() {
	for y := int32(-ball.rad); y < ball.rad; y++ {
		for x := int32(-ball.rad); x < ball.rad; x++ {
			if x*x+y*y < ball.rad*ball.rad {
				putPixel(int32(ball.pos.x)+x, int32(ball.pos.y)+y, ball.col)
			}
		}
	}
}

type paddle struct {
	pos  vec2
	w, h int32
	col  col
}

func (paddle *paddle) draw() {
	startX := int32(paddle.pos.x) - paddle.w/2
	startY := int32(paddle.pos.y) - paddle.h/2

	for y := int32(0); y < paddle.h; y++ {
		for x := int32(0); x < paddle.w; x++ {
			putPixel(startX+x, startY+y, paddle.col)
		}
	}
}

/*
Only = is the assignment operator.
:= is a part of the syntax of the Short variable declarations clause.
*/

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	ifNilPanic(err)
	defer sdl.Quit()

	// Create a window
	var whatever int32 = sdl.WINDOWPOS_UNDEFINED
	title := "Andre's game of Pong"
	win, err := sdl.CreateWindow(title, whatever, whatever, int32(w), int32(h), sdl.WINDOW_SHOWN)
	ifNilPanic(err)
	defer win.Destroy()

	// Create a renderer
	var useGPU uint32 = sdl.RENDERER_ACCELERATED
	ren, err := sdl.CreateRenderer(win, -1, useGPU)
	ifNilPanic(err)
	defer ren.Destroy()

	// Create a texture
	tex, err := ren.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	ifNilPanic(err)
	defer tex.Destroy()

	leftPaddle := paddle{vec2{100, 100}, 20, 100, col{255, 0, 0}}
	ball := ball{vec2{300., 300.}, vec2{0., 0.}, 20, col{123, 123, 0}}

	// ----------------------------- GAME LOOP -------------------------

	for {
		// Handle closing the window
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			// Click the close button
			case *sdl.QuitEvent:
				return
			}
		}

		leftPaddle.draw()
		ball.draw()

		sdl.PumpEvents()
		// mx, my, _ /*state*/ := sdl.GetMouseState()
		// fmt.Println(mx, my)

		// Update updates the given texture rectangle with new pixel data.
		tex.Update(nil, pixels, int(w*4))

		// Copy copies all of the texture to the rendering target.
		ren.Copy(tex, nil, nil)

		// Present updates the screen with any rendering performed since the previous call.
		ren.Present()

		// Hold on there boy!
		sdl.Delay(16)
	}

}
