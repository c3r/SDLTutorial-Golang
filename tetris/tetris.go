package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type color struct {
	r, g, b, a uint8
}

func drawPoint(x, y int32, color *color, renderer *sdl.Renderer) {
	renderer.SetDrawColor(color.r, color.g, color.a, 0xFF)
	renderer.DrawPoint(x, y)
}

func main() {

	sdl.Init(sdl.INIT_VIDEO)

	window, renderer, err := sdl.CreateWindowAndRenderer(1024, 768, 0)
	if err != nil {
		return
	}

	ticksPerSecond := 30
	tickInterval := uint32(1000 / ticksPerSecond)
	lastUpdateTime := uint32(0)

	initX := 0
	initY := 0

	for {
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()
		event := sdl.PollEvent()
		if event != nil && event.GetType() == sdl.QUIT {
			break
		}

		currentTime := sdl.GetTicks()
		dt := currentTime - lastUpdateTime
		timeToSleep := int32(tickInterval - dt)

		if timeToSleep > 0 {
			sdl.Delay(uint32(timeToSleep))
		}

		initX += 1
		initY += 1

		for x := initX; x < initX+250; x++ {
			for y := initY; y < initY+100; y++ {
				r := uint8(rand.Int() % 255)
				g := uint8(rand.Int() % 255)
				b := uint8(rand.Int() % 255)
				drawPoint(int32(x), int32(y), &color{r, g, b, 255}, renderer)
			}
		}

		renderer.Present()
		lastUpdateTime = currentTime
	}

	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

	// SDL_DestroyRenderer(renderer);
	// SDL_DestroyWindow(window);
	// SDL_Quit();
	// return EXIT_SUCCESS;

}
