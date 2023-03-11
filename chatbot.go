package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "image/color"
)

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

func colors() map[string]sdl.Color {
    return map[string]sdl.Color { "WHITE": { 0xFF, 0xFF, 0xFF, 0xFF }, "BLACK": { 0x00, 0x00, 0x00, 0xFF } }
}

func palette() map[string]sdl.Color {
    return map[string]sdl.Color { "BACKGROUND": colors()["WHITE"], "RECT_BG": colors()["BLACK"] }
}

// ---------------------------------------------------------------------------------------------------
// Draw realtime
func drawPoint(x, y int32, color sdl.Color, renderer *sdl.Renderer) {
    renderer.SetDrawColor(color.R, color.G, color.B, color.A)
    renderer.DrawPoint(int32(x), int32(y))
}

func drawRect(x, y int32, w, h int32, color sdl.Color, renderer *sdl.Renderer) {
    xEnd := x + w
    yEnd := y + h
    for xi := x; xi < xEnd; xi++ {
	for yi := y; yi < yEnd; yi++ {
	    drawPoint(xi, yi, color, renderer)
	}
    }
}

func clearScr(surface *sdl.Surface) {
    bgColor := palette()["BACKGROUND"]
    mappedColor := sdl.MapRGBA(surface.Format, bgColor.R, bgColor.G, bgColor.B, bgColor.A)
    surface.FillRect(&sdl.Rect { 0, 0, WINDOW_W, WINDOW_H }, mappedColor)

}
// ---------------------------------------------------------------------------------------------------
// Update surface
func updateSurfacePoint(x, y int32, clr sdl.Color, surface *sdl.Surface) {
    //new_y := y * surface.Pitch
    //new_x := x * int32(surface.Format.BytesPerPixel)
    c := color.RGBA { clr.R, clr.G, clr.B, clr.A }
    surface.Set(int(x), int(y), c)
}

func updateSurfaceRect(x, y int32, w, h int32, color sdl.Color, surface *sdl.Surface) {
    xEnd := x + w
    yEnd := y + h
    for xi := x; xi < xEnd; xi++ {
	for yi := y; yi < yEnd; yi++ {
	    updateSurfacePoint(xi, yi, color, surface)
	}
    }
}
// ---------------------------------------------------------------------------------------------------

// Handle mouse events
func handleMouseMotionEvent(state *state, evt *sdl.Event) {
    x, y, _ := sdl.GetMouseState()
    state.mouseRectangle.x = x
    state.mouseRectangle.y = y
}

func handleMouseButtonEvent(state *state, evt *sdl.Event) {
    x, y, _ := sdl.GetMouseState()
    state.visibleRectangles = append(state.visibleRectangles, rectangle{ x, y })
}

// ---------------------------------------------------------------------------------------------------

func render(state *state, renderer *sdl.Renderer, surface *sdl.Surface) {
    // Draw mouse rectangle
    drawRect(state.mouseRectangle.x, state.mouseRectangle.y, 100, 100, palette()["RECT_BG"], renderer) 
    // Udpate surface
    surface.Lock()
    for _, rectangle := range state.visibleRectangles {
    	updateSurfaceRect(rectangle.x, rectangle.y, 100, 100, palette()["RECT_BG"], surface)
    }
    surface.Unlock()
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

    window, err := sdl.CreateWindow("test", 100, 100, WINDOW_W, WINDOW_H, sdl.WINDOW_SHOWN)
    if err != nil {
	panic(err)
    }

    defer window.Destroy()

    surface, err := window.GetSurface()
    if err != nil {
	panic(err)
    }

    renderer, err := window.GetRenderer()
    if err != nil {
    	panic(err)
    }

    rect := rectangle { int32(0), int32(0) }
    state := state { true, uint32(0), rect,[]rectangle{}  }
    clearScr(surface)
    for state.running {
    	update(&state)
    	window.UpdateSurface()
    	render(&state, renderer, surface)
    	sdl.Delay(20)
    }

}
