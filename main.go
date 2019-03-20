package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type worekDanych struct {
	ScreenWidth         int32
	ScreenHeight        int32
	TableHeight         int
	TableWidth          int
	TableMarginLeft     int
	TableMarginTop      int
	ServeDirectionRight int
	ServeDirectionLeft  int
}

var (
	globals = worekDanych{
		ScreenWidth:         640,
		ScreenHeight:        480,
		TableHeight:         380,
		TableWidth:          630,
		TableMarginLeft:     5,
		TableMarginTop:      70,
		ServeDirectionRight: 1,
		ServeDirectionLeft:  -1}
)

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

var gQuit bool

func initGame() bool {

	if sdl.Init(sdl.INIT_VIDEO) != nil {
		fmt.Printf("SDL could not initialize! SDL Error: %s\n", sdl.GetError())
		return false
	}

	// Set texture filtering to linear
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Println("Warning: Linear texture filtering not enabled!")
	}

	gWindow, err := sdl.CreateWindow("SDL Tutorial GOLANG REWRITE",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		globals.ScreenWidth,
		globals.ScreenHeight,
		sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Printf("Window could not be created! SDL Error: %s\n", sdl.GetError())
		return false
	}

	_, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Printf("Renderer could not be created! SDL Error: %s\n", sdl.GetError())
		return false
	}

	// g_lpScore = new LTexture(g_renderer);
	// g_rpScore = new LTexture(g_renderer);

	// Initialize PNG loading
	var flags = img.INIT_JPG | img.INIT_PNG
	var initted = img.Init(flags)
	if (initted & flags) != flags {
		fmt.Printf("SDL_image could not initialize! SDL_image Error: %s\n", img.GetError())
		return false
	}

	return true
}

func close() {
	gWindow.Destroy()
	gRenderer.Destroy()
	img.Quit()
	sdl.Quit()
}

func main() {
	if !initGame() {
		fmt.Println("Failed to initialize!\n")
	}

	fmt.Println("Initialized sucessfully!\n")
}
