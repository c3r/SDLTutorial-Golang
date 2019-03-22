package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type consts struct {
	ScreenWidth,
	ScreenHeight,
	TableHeight,
	TableWidth,
	TableMarginLeft,
	TableMarginTop,
	ServeDirectionRight,
	ServeDirectionLeft int32
}

type velocity struct{ X, Y float64 }

type texture struct {
	renderer      *sdl.Renderer
	texture       *sdl.Texture
	width, height int32
}

type ball struct {
	sdl.Point
	stickingPaddle *paddle
	collider       sdl.Rect
	velocity       velocity
	tableRect      *sdl.Rect
}

type paddle struct {
	sdl.Point
	stickingBall           *ball
	collider               sdl.Rect
	velocity               velocity
	points, serveDirection int32
	tableRect              *sdl.Rect
}

var (
	globals = consts{
		ScreenWidth:         640,
		ScreenHeight:        480,
		TableHeight:         380,
		TableWidth:          630,
		TableMarginLeft:     5,
		TableMarginTop:      70,
		ServeDirectionRight: 1,
		ServeDirectionLeft:  -1}
)

const (
	paddleWidth    = 15
	paddleHeight   = 100
	paddleVelocity = 3
	ballWidth      = 20
	ballHeight     = 20
	ballVelocity   = 4
)

var gWindow *sdl.Window
var gRenderer *sdl.Renderer
var gLeftPadScore *texture
var gRightPadScore *texture

var gQuit bool

func initGame() bool {

	if sdl.Init(sdl.INIT_VIDEO) != nil {
		fmt.Printf("SDL could not initialize! SDL Error: %s\n", sdl.GetError())
		return false
	}

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

	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Printf("Renderer could not be created! SDL Error: %s\n", sdl.GetError())
		return false
	}

	// gLeftPadScore := texture{gRenderer, nil, 0, 0}
	// gRightPadScore := texture{gRenderer, nil, 0, 0}

	flags := img.INIT_JPG | img.INIT_PNG
	initted := img.Init(flags)
	if (initted & flags) != flags {
		fmt.Printf("SDL_image could not initialize! SDL_image Error: %s\n", img.GetError())
		return false
	}

	return true
}

func (ball *ball) handleCollision(paddle *paddle) {
	// If ball and paddle are coming in the same direction, return.
	// This prevents of multiple collisions when the paddle is coming
	// the same way that the ball is after collision.
	if float64(paddle.serveDirection)*ball.velocity.X > 0 {
		return
	}

	ball.velocity.X = math.Abs(ball.velocity.X)
	ball.velocity.X = float64(paddle.serveDirection)*ball.velocity.X + paddle.velocity.X/1.5
	ball.velocity.Y = ball.velocity.Y - paddle.velocity.Y/1.25
}

func (ball *ball) colliding(paddle *paddle) bool {
	if paddle == nil {
		return false
	}

	pc := paddle.collider
	bc := ball.collider

	lb, rb, tb, bb := bc.X, bc.X+bc.W, bc.Y, bc.Y+bc.H
	lp, rp, tp, bp := pc.X, pc.X+pc.W, pc.Y, pc.Y+pc.H

	return !(bb <= tp || tb >= bp || rb <= lp || lb >= rp)
}

func (ball *ball) move(leftPaddle, rightPaddle *paddle) {
	// if(stickToPaddle()) return
	// Move
	ball.X += int32(ball.velocity.X)
	ball.Y += int32(ball.velocity.Y)

	ball.collider.X = ball.X
	ball.collider.Y = ball.Y

	// Check for score
	// if checkForScore(lp, rp) {
	// 	return
	// }

	// Check for collisions
	if ball.colliding(leftPaddle) {
		ball.handleCollision(leftPaddle)
		return
	}

	if ball.colliding(rightPaddle) {
		ball.handleCollision(rightPaddle)
		return
	}
}

func (ball *ball) render(renderer *sdl.Renderer) {
	radius := ballWidth / 2
	renderer.SetDrawColor(0xEE, 0xE0, 0x93, 0xFF)

	center := sdl.Point{
		X: ball.X + int32(radius),
		Y: ball.Y + int32(radius)}

	for w := 0; w < radius*2; w++ {
		for h := 0; h < radius*2; h++ {
			dx := radius - w
			dy := radius - h
			if (dx*dx)+(dy*dy) <= (radius * radius) {
				renderer.DrawPoint(center.X+int32(dx), center.Y+int32(dy))
			}
		}
	}
}

func (paddle *paddle) move() {
	paddle.X += int32(paddle.velocity.X)
	paddle.collider.X = paddle.X

	mid := (paddle.tableRect.X + paddle.tableRect.W) / 2
	inMiddle := paddle.X < (mid-paddleWidth-ballWidth) || paddle.X > (mid+ballWidth)

	if (paddle.X < paddle.tableRect.X) || !inMiddle || paddle.X+paddleWidth > paddle.tableRect.X+paddle.tableRect.W {
		paddle.X -= int32(paddle.velocity.X)
		paddle.collider.X = paddle.X
	}

	paddle.Y += int32(paddle.velocity.Y)
	paddle.collider.Y = paddle.Y

	if (paddle.Y <= paddle.tableRect.Y) || (paddle.Y+paddleHeight > paddle.tableRect.Y+paddle.tableRect.H) {
		paddle.Y -= int32(paddle.Y)
		paddle.collider.Y = paddle.Y
	}
}

func (paddle *paddle) render(renderer *sdl.Renderer) {
	a := int32(4)
	r1 := sdl.Rect{
		X: paddle.collider.X,
		Y: paddle.collider.Y,
		W: paddle.collider.W - a,
		H: paddle.collider.H}

	r2 := sdl.Rect{
		X: paddle.collider.X + a,
		Y: paddle.collider.Y,
		W: paddle.collider.W - 2*a,
		H: paddle.collider.H}

	r3 := sdl.Rect{
		X: paddle.collider.X + paddle.collider.W - a,
		Y: paddle.collider.Y,
		W: a,
		H: paddle.collider.H}

	renderer.SetDrawColor(0xFF, 0x00, 0x00, 0xAA)
	renderer.FillRect(&r1)

	renderer.SetDrawColor(0xCC, 0xA7, 0x53, 0xAA)
	renderer.FillRect(&r2)

	renderer.SetDrawColor(0xFF, 0x00, 0x00, 0xAA)
	renderer.FillRect(&r3)
}

func render(ball *ball, leftPaddle, rightPaddle *paddle, renderer *sdl.Renderer) {
	renderer.SetDrawColor(0x10, 0x30, 0x60, 0xFF)
	renderer.Clear()

	leftPaddle.render(renderer)
	rightPaddle.render(renderer)
	ball.render(renderer)

	renderer.Present()
}

func update(ball *ball, leftPaddle, rightPaddle *paddle) {
	ball.move(leftPaddle, rightPaddle)
	leftPaddle.move()
	rightPaddle.move()
}

func (texture *texture) renderTexture(x, y int32, clip *sdl.Rect, angle float64, center *sdl.Point, flip sdl.RendererFlip) {
	renderQuad := sdl.Rect{
		X: x,
		Y: y,
		W: texture.width,
		H: texture.height}

	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	texture.renderer.CopyEx(texture.texture, clip, &renderQuad, angle, center, flip)
}

func (texture *texture) free() {
	if texture != nil {
		texture.texture.Destroy()
		texture.texture = nil
		texture.height = 0
		texture.width = 0
	}
}

func close() {
	gWindow.Destroy()
	gRenderer.Destroy()
	img.Quit()
	sdl.Quit()
}

func loadTextureFromFile(texture *texture, path string) bool {
	loadedSurface, err := img.Load(path)
	if err != nil {
		fmt.Printf("Unable to load image %s! SDL_image Error: %s\n", path, img.GetError())
		return false
	}

	loadedSurface.SetColorKey(true, sdl.MapRGB(loadedSurface.Format, 0, 0xFF, 0xFF))

	newTexture, err := texture.renderer.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		fmt.Printf("Unable to create texture from %s! SDL Error: %s\n", path, sdl.GetError())
		return false
	}

	texture.width = loadedSurface.W
	texture.height = loadedSurface.H

	loadedSurface.Free()

	texture.texture = newTexture
	return texture.texture != newTexture
}

func main() {
	if !initGame() {
		fmt.Println("Failed to initialize!")
	}

	tableRect := sdl.Rect{
		X: globals.TableMarginLeft,
		Y: globals.TableMarginTop,
		W: globals.TableWidth,
		H: globals.TableHeight}

	leftPaddle := paddle{}
	leftPaddle.X = paddleWidth
	leftPaddle.Y = globals.ScreenHeight/2 - paddleHeight/2
	leftPaddle.collider = sdl.Rect{
		X: leftPaddle.X,
		Y: leftPaddle.Y,
		W: paddleWidth,
		H: paddleHeight}

	leftPaddle.points = 0
	leftPaddle.serveDirection = 1
	leftPaddle.stickingBall = nil
	leftPaddle.velocity = velocity{0, 0}
	leftPaddle.tableRect = &tableRect

	rightPaddle := paddle{}
	rightPaddle.X = globals.ScreenWidth - 2*paddleWidth
	rightPaddle.Y = globals.ScreenHeight/2 - paddleHeight/2
	rightPaddle.collider = sdl.Rect{
		X: rightPaddle.X,
		Y: rightPaddle.Y,
		W: paddleWidth,
		H: paddleHeight}
	rightPaddle.points = 0
	rightPaddle.serveDirection = -1
	rightPaddle.stickingBall = nil
	rightPaddle.velocity = velocity{0, 0}
	rightPaddle.tableRect = &tableRect

	ball := ball{}
	ball.X = globals.ScreenWidth/2 - ballWidth/2
	ball.Y = globals.ScreenHeight/2 - ballHeight/2
	ball.collider = sdl.Rect{
		X: ball.X,
		Y: ball.Y,
		W: ballWidth,
		H: ballHeight}
	ball.stickingPaddle = nil
	ball.velocity = velocity{2, 0}
	ball.tableRect = &tableRect

	fmt.Println("Initialized sucessfully!")

	ticksPerSecond := 60
	tickInterval := uint32(1000 / ticksPerSecond)
	lastUpdateTime := uint32(0)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		currentTime := sdl.GetTicks()
		dt := currentTime - lastUpdateTime
		timeToSleep := int32(tickInterval - dt)

		if timeToSleep > 0 {
			sdl.Delay(uint32(timeToSleep))
		}

		update(&ball, &leftPaddle, &rightPaddle)
		render(&ball, &leftPaddle, &rightPaddle, gRenderer)

		lastUpdateTime = currentTime
	}
}
