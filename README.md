# DISCONTINUED
# OverlayGo
A lightweight library on Go that allows you to work with OpenGl with a transparent overlay.
![image](https://github.com/MrRoblick/OverlayGo/assets/61147166/ae3ff41f-f400-463e-8c8d-0600128fba3a)

Example
```go
package main
import(
	"DrawerGO/Overlay"
)
func main() {
	renderer := Overlay.New()
	renderer.Run()
	font := renderer.LoadFont("Fonts/Vcr.ttf", 48)
	imageId, w, h := renderer.LoadImage("Resources/MinosPrime.png")
	defer renderer.Dispose()
	for renderer.IsOpen() {
		x, y := renderer.GetMousePosition()
		sX, sY := renderer.GetMonitorSize()
		renderer.AnchorPoint(0.5, 0.5)
		renderer.SetColor(0, 255, 0, 255)
		renderer.DrawText(x, y, 30, 0, 0, font, .8, `Hello World!`)
		renderer.SetColor(255, 0, 0, 255)
		renderer.DrawOutlineRect(500, 500, 100, 100)
		renderer.SetColor(0, 0, 255, 255)
		renderer.DrawLine(600, 0, 100, 100)
		renderer.DrawPolygon(0, 0,
			0, 50,
			50, 50,
		)
		renderer.SetColor(255, 255, 0, 255)
		renderer.AnchorPoint(0, 0)
		renderer.DrawCircle(0, float32(sY)/2, 120, 120)
		renderer.ResetColor()
		renderer.RotateByDeg(20)
		renderer.DrawImage(float32(sX)/2, float32(sY)/2, float32(w)*0.2, float32(h)*0.2, imageId)
		renderer.Render()
	}
}
```
