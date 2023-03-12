package main

import (
    "github.com/veandco/go-sdl2/sdl"
//    "image/color"
    "math"
    "math/rand"
    "sort"
    "fmt"
//    "io"
//    "os"
)

const (
    PIXELSIZE = 3
    GRIDSIZE = 30 * PIXELSIZE
    GRIDS_HORIZ = 30
    GRIDS_VERTI = 20
    WINDOW_W = GRIDSIZE * GRIDS_HORIZ
    WINDOW_H = GRIDSIZE * GRIDS_VERTI
    RECT_W = GRIDSIZE * 5
    RECT_H = GRIDSIZE * 5
)

// ===============================================================================================
// State handling code
// ===============================================================================================

var GLOBAL_ID int

func nextId() int {
    GLOBAL_ID++
    return GLOBAL_ID 
}

// DiagramCtx class
type DiagramCtx struct {
    running bool
    elements map[int]DiagramElement
}

type DiagramElement struct {
    id int
}

func (this *DiagramCtx) createDiagramElement(id int) DiagramElement {
    this.elements[id] = DiagramElement { id } 
    return this.elements[id]
}

func (this *DiagramCtx) isRunning() bool {
    return this.running
}

func (this *DiagramCtx) stop() {
    this.running = false
}

func update(diagram *DiagramCtx) {
}

// ===============================================================================================
// Drawing code
// ===============================================================================================
type RenderFunc func(*Graphics, *DrawableShape) *sdl.Texture

type DrawableShape struct {
    rect *sdl.Rect
    texture *sdl.Texture
    layer int
    rendered bool
}

func (this *DrawableShape) log() {
    fmt.Printf("(%d, %d): %d x %d \n", this.rect.X, this.rect.Y, this.rect.W, this.rect.H)
}

func (this *DrawableShape) getLayer() int {
    return this.layer
}

func (this *DrawableShape) isRendered() bool {
    return this.rendered
}

func (this *DrawableShape) setX(x int32) {
    this.rect.X = snapToGrid(x)
}

func (this *DrawableShape) setY(y int32) {
    this.rect.Y = snapToGrid(y)
}

func snapToGrid(value int32) int32 {
    return int32( math.Round( float64(value) / GRIDSIZE ) * GRIDSIZE )
}

func (this *DrawableShape) getTexture() *sdl.Texture {
    return this.texture
}

func (this *DrawableShape) getRectangle() *sdl.Rect {
    return this.rect
}

type Graphics struct {
    render *sdl.Renderer
    shapes map[int]DrawableShape
    currentLayer int
    colorRotationArray []sdl.Color
}

func (this *Graphics) getCurrentLayer() int {
    return this.currentLayer
}

func (this *Graphics) setCurrentLayer(layer int) {
    this.currentLayer = layer
}

func (this *Graphics) draw() {
    sortedIds := make([]int, 0, len(this.shapes))
    for id := range this.shapes {
	sortedIds = append(sortedIds, id)
    }
    sort.Ints(sortedIds)

    layer := 0 
    renderedShapes := 0
    for renderedShapes < len(this.shapes) {
	for _, id := range sortedIds {
	    shape, ok := this.shapes[id]
	    if ok && shape.getLayer() == layer && !shape.isRendered() {
		this.render.Copy(shape.getTexture(), nil, shape.getRectangle())	
		shape.rendered = true
		renderedShapes++
	    }
	}
	layer++
    }

    for id := range sortedIds {
	shape := this.shapes[id] 
	shape.rendered = false
    }
    
    this.render.Present()
}

func (this *Graphics) createDrawableShape(id, w, h int, renderFunc RenderFunc) *DrawableShape {
    rect := sdl.Rect { 0, 0, int32(w), int32(h) }
    texture, _ := this.render.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, int32(w), int32(h))
    shape := DrawableShape { &rect, texture, this.getCurrentLayer(), false }    
    texture = renderFunc(this, &shape)
    this.shapes[id] = shape
    return &shape
}

// ===============================================================================================
// Rendering code
// ===============================================================================================
func renderRectTexture(graphics *Graphics, shape *DrawableShape) *sdl.Texture {
    print("Rendering rect...\n")
    texture := shape.getTexture()
    graphics.render.SetRenderTarget(texture)
    rect := shape.getRectangle()
    // Border
  
    color := graphics.colorRotationArray[rand.Intn(len(graphics.colorRotationArray))]
    graphics.render.SetDrawColor(color.R - 30, color.G - 30, color.B - 30, color.A)
    graphics.render.Clear()
    // Inner Color
    graphics.render.SetDrawColor(color.R, color.G, color.B, color.A)
    p := int32( PIXELSIZE )
    innerRect := sdl.Rect { 
        rect.X + p, 
        rect.Y + p, 
        rect.W - 2*p, 
        rect.H - 2*p,
    }
    graphics.render.FillRect(&innerRect)
    graphics.render.SetRenderTarget(nil)
    return texture
}

func renderBgTexture(graphics *Graphics, shape *DrawableShape) (*sdl.Texture) {
    print("Rendering background...\n")
    texture := shape.getTexture()
    graphics.render.SetRenderTarget(texture)
    graphics.render.SetDrawColor(0xF1, 0xF1, 0xF1, 0xFF)
    graphics.render.Clear()
    graphics.render.SetDrawColor(0xDD, 0xDD, 0xDD, 255)
    for i := GRIDSIZE; i < WINDOW_W; i += GRIDSIZE {
	graphics.render.DrawLine(int32(i), 0, int32(i), WINDOW_H)
    }
    for i := GRIDSIZE; i < WINDOW_H; i += GRIDSIZE {
	graphics.render.DrawLine(0, int32(i), WINDOW_W, int32(i))
    }
    graphics.render.SetRenderTarget(nil)
    return texture
}

func renderMenuButton(graphics *Graphics, shape *DrawableShape) (*sdl.Texture) {
    print("Rendering menu button...\n")
    texture := shape.getTexture()
    graphics.render.SetRenderTarget(texture)
    // Shadow
    graphics.render.SetDrawColor(0x33, 0x33, 0x33, 0xFF)
    graphics.render.Clear()
    // White
    rect := shape.getRectangle()
    p := int32( PIXELSIZE )
    innerRect := sdl.Rect { 
        rect.X, 
        rect.Y, 
        rect.W - p, 
        rect.H - p,
    }
    graphics.render.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
    graphics.render.FillRect(&innerRect)
    // Dark gray
    innerRect = sdl.Rect { 
        rect.X + p, 
        rect.Y + p, 
        rect.W - 2*p, 
        rect.H - 2*p,
    }
    graphics.render.SetDrawColor(0x66, 0x66, 0x66, 0xFF)
    graphics.render.FillRect(&innerRect)
    // Light gray
    innerRect = sdl.Rect { 
        rect.X + p, 
        rect.Y + p, 
        rect.W - 3*p, 
        rect.H - 3*p,
    }
    graphics.render.SetDrawColor(0xCC, 0xCC, 0xCC, 0xFF)
    graphics.render.FillRect(&innerRect)
    // Front color
    innerRect = sdl.Rect { 
        rect.X + 2*p, 
        rect.Y + 2*p, 
        rect.W - 4*p, 
        rect.H - 4*p,
    }
    graphics.render.SetDrawColor(0xBB, 0xBB, 0xBB, 0xFF)
    graphics.render.FillRect(&innerRect)

    graphics.render.SetRenderTarget(nil)
    return texture
}

// =========================================================================================================
// Main code
// =========================================================================================================
func main() {
    // Setup
    sdl.Init(sdl.INIT_EVERYTHING) 
    defer sdl.Quit()
    window, _ := sdl.CreateWindow("test", 100, 100, WINDOW_W, WINDOW_H, sdl.WINDOW_SHOWN)
    defer window.Destroy()
    renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

    colorRotationArray := []sdl.Color{
	sdl.Color { 0x89, 0xbd, 0x9e, 0xFF },
	sdl.Color { 0xf0, 0xc9, 0x87, 0xFF },
	sdl.Color { 0xdb, 0x4c, 0x40, 0xFF },
    }

    graphics := Graphics { renderer, map[int]DrawableShape{}, 1, colorRotationArray }
    diagram := DiagramCtx { true, map[int]DiagramElement{} }
    graphics.createDrawableShape(nextId(), WINDOW_W, WINDOW_H, renderBgTexture) 

    // Create only graphics for background, we don't need it as diagram element
    for i:=0; i<8; i++ {
	menuButton := graphics.createDrawableShape(nextId(), GRIDSIZE, GRIDSIZE, renderMenuButton)
	menuButton.setX(0)
	menuButton.setY(GRIDSIZE + menuButton.getRectangle().H * int32(i))
    }

    // Create the mouse rectangle
    //id := 100
    //diagram.createDiagramElement(id)
    //mouseRectangle := graphics.createDrawableShape(id, RECT_W, RECT_H, renderRectTexture) 

    for diagram.isRunning() {
	graphics.draw()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	    switch event.(type) {
		case *sdl.MouseMotionEvent:
		    //x, y, _ := sdl.GetMouseState()
		    //mouseRectangle.setX(x)
		    //mouseRectangle.setY(y)
		case *sdl.MouseButtonEvent:
		    //x, y, _ := sdl.GetMouseState()
		    //if event.GetType() == sdl.MOUSEBUTTONDOWN {
		    //    id := nextId()
		    //    diagram.createDiagramElement(id)
		    //    shape := graphics.createDrawableShape(id, RECT_W, RECT_H, renderRectTexture) 
		    //    shape.setX(x)
		    //    shape.setY(y)
		    //}
		case *sdl.QuitEvent:
		    println("Quit")
		    diagram.stop()
		    break
	    }
	}
	sdl.Delay(20)
    }

}
