package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type color struct {
	r, g, b, a uint8
}

type renderingCtx struct {
	fps            uint32
	lastUpdateTime uint32
	currTime       uint32
}

const (
	gSideLen          = 10
	gTetrominoRectNum = 4
	gFPS              = 10
	gFastFPS          = 60
	gWindowWidth      = 1024
	gWindowHeight     = 768
)

func drawTetromino(tetromino [16]int, x, y int, position int, renderer *sdl.Renderer) {
	yi := y
	xi := x
	for i := 0; i < gTetrominoRectNum; i++ {
		yi = y + ((tetromino[i+(position*4)] / gTetrominoRectNum) * gSideLen)
		xi = x + ((tetromino[i+(position*4)] % gTetrominoRectNum) * gSideLen)
		drawSquare(xi, yi, gSideLen, &color{255, 255, 255, 255}, renderer)
	}
}

func drawPoint(x, y int, color *color, renderer *sdl.Renderer) {
	renderer.SetDrawColor(color.r, color.g, color.a, 0xFF)
	renderer.DrawPoint(int32(x), int32(y))
}

func drawSquare(x, y int, a int, color *color, renderer *sdl.Renderer) {
	drawRect(x, y, a, a, color, renderer)
}

func drawRect(x, y int, w, h int, color *color, renderer *sdl.Renderer) {
	xEnd := x + w
	yEnd := y + h
	for xi := x; xi < xEnd; xi++ {
		for yi := y; yi < yEnd; yi++ {
			drawPoint(xi, yi, color, renderer)
		}
	}
}

func (ctx *renderingCtx) sleep(keyState []uint8) {
	tts := ctx.timeToSleep(keyState)
	if tts > 0 {
		sdl.Delay(uint32(tts))
	}
	ctx.lastUpdateTime = ctx.currTime
}

func (ctx *renderingCtx) timeToSleep(keyState []uint8) int32 {
	ctx.currTime = sdl.GetTicks()
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		ctx.fps = gFastFPS
	} else {
		ctx.fps = gFPS
	}
	return int32(uint32(1000/ctx.fps) - (ctx.currTime - ctx.lastUpdateTime))
}

func main() {

	t1 := [16]int{1, 5, 8, 9, 0, 4, 5, 6, 1, 2, 5, 9, 0, 1, 2, 6}
	t2 := [16]int{1, 5, 9, 10, 0, 1, 2, 4, 0, 1, 5, 9, 2, 4, 5, 6}
	t3 := [16]int{0, 4, 8, 12, 0, 1, 2, 3, 0, 4, 8, 12, 0, 1, 2, 3}
	t4 := [16]int{1, 4, 5, 6, 1, 5, 6, 9, 0, 1, 2, 5, 1, 4, 5, 9}
	t5 := [16]int{0, 1, 4, 5, 0, 1, 4, 5, 0, 1, 4, 5, 0, 1, 4, 5}

	sdl.Init(sdl.INIT_VIDEO)

	window, renderer, err := sdl.CreateWindowAndRenderer(gWindowWidth, gWindowHeight, 0)
	if err != nil {
		return
	}

	keyState := sdl.GetKeyboardState()
	renderingCtx := renderingCtx{uint32(60), uint32(0), uint32(0)}
	position := 0
	ypositionOffset := 0

	for {
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			if event.GetType() == sdl.KEYDOWN {
				position++
			}

			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderingCtx.sleep(keyState)

		drawTetromino(t1, 100, ypositionOffset, position%4, renderer)
		drawTetromino(t2, 200, ypositionOffset, position%4, renderer)
		drawTetromino(t3, 300, ypositionOffset, position%4, renderer)
		drawTetromino(t4, 400, ypositionOffset, position%4, renderer)
		drawTetromino(t5, 500, ypositionOffset, position%4, renderer)

		ypositionOffset += gSideLen

		renderer.Present()
	}

	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

}
