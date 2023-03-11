package main

import "github.com/veandco/go-sdl2/sdl"

type color struct {
    r, g, b, a uint8
}

const (
    WINDOW_W = 1920
    WINDOW_H = 1080
)

type rectangle struct {
    x, y int32
}

type state struct {
    running bool
    runningTime uint32
    mouseRectangle rectangle 
    visibleRectangles []rectangle
}

func colors() map[string]color {
    return map[string]color { "WHITE": { 0xFF, 0xFF, 0xFF, 0xFF }, "BLACK": { 0x00, 0x00, 0x00, 0xFF } }
}

func palette() map[string]color {
    return map[string]color { "BACKGROUND": colors()["WHITE"], "RECT_BG": colors()["BLACK"] }
}

func drawPoint(x, y int32, color color, renderer *sdl.Renderer) {
    renderer.SetDrawColor(color.r, color.g, color.b, 0xFF)
    renderer.DrawPoint(int32(x), int32(y))
}

func drawRect(x, y int32, w, h int32, color color, renderer *sdl.Renderer) {
    xEnd := x + w
    yEnd := y + h
    for xi := x; xi < xEnd; xi++ {
	for yi := y; yi < yEnd; yi++ {
	    drawPoint(xi, yi, color, renderer)
	}
    }
}

func clearScr(renderer *sdl.Renderer) {
    bgColor := palette()["BACKGROUND"]
    renderer.SetDrawColor(bgColor.r, bgColor.g, bgColor.b, 0xFF)
    renderer.Clear()
}

func handleMouseMotionEvent(state *state, evt *sdl.Event) {
    x, y, _ := sdl.GetMouseState()
    state.mouseRectangle.x = x
    state.mouseRectangle.y = y
}

func handleMouseButtonEvent(state *state, evt *sdl.Event) {
    x, y, _ := sdl.GetMouseState()
    state.visibleRectangles = append(state.visibleRectangles, rectangle{ x, y })
}

func render(state *state, renderer *sdl.Renderer) {
    // Draw mouse rectangle
    drawRect(state.mouseRectangle.x, state.mouseRectangle.y, 100, 100, palette()["RECT_BG"], renderer) 
    for _, rectangle := range state.visibleRectangles {
    	drawRect(rectangle.x, rectangle.y, 100, 100, palette()["RECT_BG"], renderer)
    }
    renderer.Present()	
}

func update(state *state) {
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	switch event.(type) {
	    case *sdl.MouseMotionEvent:
		handleMouseMotionEvent(state, &event)	
	    case *sdl.MouseButtonEvent:
	    	handleMouseButtonEvent(state, &event)
	    case *sdl.QuitEvent:
		println("Quit")
		state.running = false
		break
	}
    }
}

func main() {
    
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	panic(err)
    }
    defer sdl.Quit()

    window, renderer, err := sdl.CreateWindowAndRenderer(WINDOW_W, WINDOW_H, 0)
    if err != nil {
	panic(err)
    }
    defer window.Destroy()

    if err != nil {
	panic(err)
    }

    rect := rectangle { int32(0), int32(0) }
    state := state { true, uint32(0), rect,[]rectangle{}  }
    for state.running {
	clearScr(renderer)
    	update(&state)
    	render(&state, renderer)
    }

}
