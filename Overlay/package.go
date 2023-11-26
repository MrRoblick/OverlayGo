package Overlay

import (
	"DrawerGO/Overlay/GlTools"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gonutz/w32/v2"
	"runtime"
)

// var loadedImages map[string][]byte
var currentColor = [4]float32{1, 1, 1, 1}
var defaultColor = [4]float32{1, 1, 1, 1}
var currentRotation float32 = 0
var currentPositionX, currentPositionY float32 = 0, 0
var currentZIndex uint32 = 0
var currentAnchorPointX, currentAnchorPointY float32 = 0, 0
var currentFill = true

type ProgressBarDirection byte

const (
	PROGRESS_BAR_DIRECTION_LEFT ProgressBarDirection = iota
	PROGRESS_BAR_DIRECTION_RIGHT
	PROGRESS_BAR_DIRECTION_TOP
	PROGRESS_BAR_DIRECTION_BOTTOM
	PROGRESS_BAR_DIRECTION_CENTER
)

//var thickness float32 = 1

func New() App {
	window := newWindow("DrawerOverlayWindow")
	initGl()
	ctx := newContext(window.glfwWindow, window.hwnd)

	ctx.Init()
	return App{
		context: ctx,
		window:  window,
		isRun:   true,
	}
}

func (app *App) Run() {
	app.isRun = true
	runtime.LockOSThread()
}
func (app *App) Stop() {
	app.isRun = false
}
func (app *App) IsOpen() bool {
	return !app.window.glfwWindow.ShouldClose() && app.isRun
}
func (app *App) Dispose() {
	app.isRun = false
	app.context.ClearAll()
	glfw.Terminate()
}
func (app *App) Render() {
	app.context.Render()
	app.context.ClearAll()
	app.ResetRotate()
	app.ResetColor()
	app.ResetAnchorPoint()
	currentZIndex = 0
}

func (app *App) GetMousePosition() (float32, float32) {
	x, y := app.window.glfwWindow.GetCursorPos()
	return float32(x), float32(y)
}

func (app *App) GetMonitorSize() (int, int) {
	x, y := app.window.glfwWindow.GetSize()
	y -= 1
	return x, y
}

/*

Renderer

*/

func (app *App) SetColor(R, G, B, A byte) {
	currentColor = [4]float32{
		float32(R) / 255.0,
		float32(G) / 255.0,
		float32(B) / 255.0,
		float32(A) / 255.0,
	}

}
func (app *App) ResetColor() {
	currentColor = defaultColor
}
func (app *App) ResetAnchorPoint() {
	currentAnchorPointX = 0
	currentAnchorPointY = 0
}
func (app *App) DrawLine(X1, Y1, X2, Y2 float32) {
	currentZIndex++
	app.context.lines = append(app.context.lines, Line{
		X1: X1,
		Y1: Y1,

		X2:           X2,
		Y2:           Y2,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}
func (app *App) DrawOutlineRect(X, Y, Width, Height float32) {
	currentZIndex++
	app.context.outlineRects = append(app.context.outlineRects, OutlineRect{
		X: X,
		Y: Y,

		Width:        Width,
		Height:       Height,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
	})
}
func (app *App) DrawRect(X, Y, Width, Height float32) {
	currentZIndex++
	app.context.rects = append(app.context.rects, Rect{
		X: X,
		Y: Y,

		Width:        Width,
		Height:       Height,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}

// DrawImage To get the ImageId, you need to upload an image via LoadImage
func (app *App) DrawImage(X, Y, Width, Height float32, ImageId uint32) {
	currentZIndex++
	app.context.images = append(app.context.images, Image{
		X:            X,
		Y:            Y,
		Width:        Width,
		Height:       Height,
		Image:        ImageId,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}
func (app *App) DrawText(X, Y, Size, Width, Height float32, Font uint32, Interval float32, text string) {
	if len(text) == 0 {
		return
	}
	currentZIndex++
	app.context.texts = append(app.context.texts, Text{
		X:            X,
		Y:            Y,
		Size:         Size,
		Text:         text,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
		Width:        Width,
		Height:       Height,
		FontTexture:  Font,
		Interval:     Interval,
	})
}
func (app *App) DrawPolygon(X1, Y1, X2, Y2, X3, Y3 float32) {
	currentZIndex++
	app.context.polygons = append(app.context.polygons, Polygon{
		X1: X1,
		Y1: Y1,

		X2: X2,
		Y2: Y2,

		X3:           X3,
		Y3:           Y3,
		Color:        currentColor,
		Rotation:     currentRotation,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}

func (app *App) DrawCircle(X, Y, ScaleX, ScaleY float32) {
	currentZIndex++
	app.context.circles = append(app.context.circles, Circle{
		X:            X,
		Y:            Y,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ScaleX:       ScaleX,
		ScaleY:       ScaleY,
		Color:        currentColor,
		Rotation:     currentRotation,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}
func clamp[T float32 | float64 | int | int32 | int64 | int8 | int16](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {

		return max
	}
	return value
}
func (app *App) DrawProgressBar(X, Y, Width, Height, Value, Max float32, Direction ProgressBarDirection) {
	var _w, _h float32 = 1, 1
	var x, y float32
	switch Direction {
	case PROGRESS_BAR_DIRECTION_RIGHT:
		{
			_w = clamp[float32](Value/Max, 0, 1)
			x = Width * (currentAnchorPointX - _w/2)
			y = Height * (currentAnchorPointY - _h/2)
		}
	case PROGRESS_BAR_DIRECTION_LEFT:
		{
			_w = -clamp[float32](Value/Max, 0, 1)
			x = Width * (-currentAnchorPointX - _w/2)
			y = Height * (currentAnchorPointY - _h/2)
		}
	case PROGRESS_BAR_DIRECTION_TOP:
		{
			_h = clamp[float32](Value/Max, 0, 1)
			x = Width * (currentAnchorPointX - _w/2)
			y = Height * (-currentAnchorPointY + _h/2)
		}
	case PROGRESS_BAR_DIRECTION_BOTTOM:
		{
			_h = -clamp[float32](Value/Max, 0, 1)
			x = Width * (currentAnchorPointX - _w/2)
			y = Height * (currentAnchorPointY - _h/2)
		}
	case PROGRESS_BAR_DIRECTION_CENTER:
		{
			_w = clamp[float32](Value/Max, 0, 1)
		}
	}
	if Direction == PROGRESS_BAR_DIRECTION_RIGHT {

	}
	currentZIndex++
	app.context.rects = append(app.context.rects, Rect{
		X: X,
		Y: Y,

		Width:        Width,
		Height:       Height,
		Color:        [4]float32{currentColor[0] / 2, currentColor[1] / 2, currentColor[2] / 2, currentColor[3] / 2},
		Rotation:     0,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})

	currentZIndex++
	app.context.rects = append(app.context.rects, Rect{
		X: X - x,
		Y: Y - y,

		Width:        _w * Width,
		Height:       _h * Height,
		Color:        currentColor,
		Rotation:     0,
		TransformX:   currentPositionX,
		TransformY:   currentPositionY,
		ZIndex:       currentZIndex,
		AnchorPointX: currentAnchorPointX,
		AnchorPointY: currentAnchorPointY,
		Fill:         currentFill,
	})
}
func (app *App) RotateByRad(Deg float32) {
	currentRotation = mgl32.DegToRad(Deg)
}
func (app *App) RotateByDeg(Deg float32) {
	currentRotation = Deg
}
func (app *App) Translate(X, Y float32) {
	currentPositionX = X
	currentPositionY = Y
}
func (app *App) GetFps() float32 {
	return app.context.GetFPS()
}
func (app *App) GetDeltaTime() float32 {
	return app.context.GetDeltaTime()
}

func (app *App) LoadImage(path string) (uint32, int, int) {
	tex := GlTools.MakeTexture(true)
	width, height, err := GlTools.UploadTexture(tex, path)
	if err != nil {
		fmt.Println(err)
		gl.DeleteTextures(1, &tex)
		return 0, 0, 0
	}
	return tex, width, height
}
func (app *App) LoadFont(path string, FontSize float32) uint32 {
	return app.context.LoadFont(path, FontSize)
}
func (app *App) DeleteImage(imgId uint32) {
	gl.DeleteTextures(1, &imgId)
}
func (app *App) AnchorPoint(X, Y float32) {
	currentAnchorPointX = X
	currentAnchorPointY = Y
}
func (app *App) GetTime() float32 {
	return float32(glfw.GetTime())
}
func (app *App) ResetRotate() {
	currentRotation = 0
}
func (app *App) IsKeyDown(Key int) bool {
	return w32.GetKeyState(int(Key)) > 1
}
func (app *App) IsMouseButton1Down() bool {
	return w32.GetAsyncKeyState(w32.VK_LBUTTON) > 1
}
func (app *App) IsMouseButton2Down() bool {
	return w32.GetAsyncKeyState(w32.VK_RBUTTON) > 1
}
func (app *App) IsShiftDown() bool {
	return w32.GetAsyncKeyState(w32.VK_SHIFT) > 1
}
func (app *App) IsCtrlDown() bool {
	return w32.GetAsyncKeyState(w32.VK_CONTROL) > 1
}
func (app *App) IsBackspaceDown() bool {
	return w32.GetAsyncKeyState(w32.VK_BACK) > 1
}
func (app *App) IsEnterDown() bool {
	return w32.GetAsyncKeyState(w32.VK_RETURN) > 1
}
func (app *App) SetFillMode() {
	currentFill = true
}
func (app *App) SetWireframeMode() {
	currentFill = false
}
func (app *App) SetMode(fill bool) {
	currentFill = fill
}
