package Overlay

import (
	"DrawerGO/Overlay/GlTools"
	"DrawerGO/Overlay/Mesh"
	"DrawerGO/Overlay/Shader"
	"DrawerGO/Overlay/Type"
	"DrawerGO/Overlay/ttf2atlas"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gonutz/w32/v2"
	"strings"
	"unsafe"
)

type vec2 struct {
	X, Y int
}
type atlasFont struct {
	Size     vec2
	FontSize float32
}

func GetDisplaySize() (int32, int32) {
	hMonitor := w32.MonitorFromPoint(0, 0, w32.MONITOR_DEFAULTTOPRIMARY)

	var monitorInfo w32.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	w32.GetMonitorInfo(hMonitor, &monitorInfo)

	width := monitorInfo.RcMonitor.Right - monitorInfo.RcMonitor.Left
	height := monitorInfo.RcMonitor.Bottom - monitorInfo.RcMonitor.Top
	return width, height
}

func initGlfw(name string) (*glfw.Window, w32.HWND) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.Samples, 2)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	glfw.WindowHint(glfw.Focused, glfw.False)
	glfw.WindowHint(glfw.FocusOnShow, glfw.False)

	w, h := GetDisplaySize()
	window, err := glfw.CreateWindow(int(w), int(h+1), name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetAttrib(glfw.Floating, glfw.True)

	window.Focus()
	hwnd := w32.GetActiveWindow()
	w32.SetWindowLong(
		hwnd,
		w32.GWL_EXSTYLE,
		w32.WS_EX_TRANSPARENT|w32.WS_EX_LAYERED|w32.WS_EX_TOOLWINDOW,
	)
	w32.SetWindowPos(
		hwnd,
		w32.HWND_TOPMOST,
		0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE,
	)
	window.Show()

	return window, hwnd
}
func initGl() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Open GL: ", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.FrontFace(gl.FRONT_FACE)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

type Line struct {
	X1, Y1, X2, Y2                                     float32
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
	Fill                                               bool
}
type Rect struct {
	X, Y, Width, Height                                float32
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
	Fill                                               bool
}
type Image struct {
	X, Y, Width, Height                                float32
	Image                                              uint32
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
	Fill                                               bool
}
type Polygon struct {
	X1, Y1, X2, Y2, X3, Y3                             float32
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
	Fill                                               bool
}
type Text struct {
	X, Y, Size                                         float32
	Text                                               string
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
	Fill                                               bool
	FontTexture                                        uint32
	Width, Height                                      float32
	Interval                                           float32
}
type Circle struct {
	X, Y, TransformX, TransformY, ScaleX, ScaleY, AnchorPointX, AnchorPointY, Rotation float32
	Color                                                                              [4]float32
	ZIndex                                                                             uint32
	Fill                                                                               bool
}
type OutlineRect struct {
	X, Y, Width, Height                                float32
	Color                                              [4]float32
	Rotation                                           float32
	TransformX, TransformY, AnchorPointX, AnchorPointY float32
	ZIndex                                             uint32
}
type Context struct {
	lines        []Line
	rects        []Rect
	images       []Image
	polygons     []Polygon
	texts        []Text
	circles      []Circle
	outlineRects []OutlineRect
	mainProgram  uint32

	lineRenderObject,
	rectRenderObject,
	imageRenderObject,
	polygonRenderObject,
	circleRenderObject, outlineRectRenderObject, textRenderObject GlTools.RenderObject

	window *glfw.Window

	colorUniform, modelUniform, cameraUniform, textureUniform, textureEnabledUniform int32
	vertexAttributeLocation, uvAttributeLocation                                     uint32
	vertexAttribute, uvAttribute                                                     GlTools.Attribute
	hwnd                                                                             w32.HWND
	lastTime, deltaTime, fps                                                         float32
	loadedFonts                                                                      map[uint32]atlasFont
}

type Window struct {
	name       string
	hwnd       w32.HWND
	glfwWindow *glfw.Window
}
type App struct {
	context Context
	window  Window
	isRun   bool
}

func newContext(window *glfw.Window, hwnd w32.HWND) Context {
	line := GlTools.NewRenderObject(Type.Line)

	rect := GlTools.NewRenderObject(Type.Rectangle)
	rect.UploadMesh(Mesh.Rect())

	circle := GlTools.NewRenderObject(Type.Circle)
	circle.UploadMesh(Mesh.Circle(360))

	polygon := GlTools.NewRenderObject(Type.Polygon)

	img := GlTools.NewRenderObject(Type.Image)
	img.UploadMesh(Mesh.Rect())

	text := GlTools.NewRenderObject(Type.Text)

	outlineRect := GlTools.NewRenderObject(Type.OutlineRect)
	outlineRect.UploadMesh(Mesh.Rect())
	return Context{
		lines:                   []Line{},
		rects:                   []Rect{},
		images:                  []Image{},
		polygons:                []Polygon{},
		circles:                 []Circle{},
		mainProgram:             0,
		lineRenderObject:        line,
		rectRenderObject:        rect,
		circleRenderObject:      circle,
		polygonRenderObject:     polygon,
		imageRenderObject:       img,
		outlineRectRenderObject: outlineRect,
		textRenderObject:        text,
		window:                  window,
		hwnd:                    hwnd,
		loadedFonts:             map[uint32]atlasFont{},
	}
}
func newWindow(name string) Window {
	win, hwnd := initGlfw(name)
	glfw.SwapInterval(0)
	return Window{
		name:       name,
		glfwWindow: win,
		hwnd:       hwnd,
	}
}

func (ctx *Context) Init() {
	vertShader, err := GlTools.CompileShader(Shader.RendererVertexShader, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := GlTools.CompileShader(Shader.RendererFragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	prog, err := GlTools.NewProgram(vertShader, fragShader)

	if err != nil {
		panic(err)
	}
	ctx.mainProgram = prog

	ctx.colorUniform = gl.GetUniformLocation(prog, gl.Str("BaseColor\x00"))
	ctx.modelUniform = gl.GetUniformLocation(prog, gl.Str("Model\x00"))
	ctx.cameraUniform = gl.GetUniformLocation(prog, gl.Str("Camera\x00"))
	ctx.textureUniform = gl.GetUniformLocation(prog, gl.Str("tex\x00"))
	ctx.textureEnabledUniform = gl.GetUniformLocation(prog, gl.Str("texEnabled\x00"))
	ctx.vertexAttributeLocation = uint32(gl.GetAttribLocation(prog, gl.Str("Vert\x00")))
	ctx.uvAttributeLocation = uint32(gl.GetAttribLocation(prog, gl.Str("Uv\x00")))

	ctx.vertexAttribute = GlTools.NewAttribute(ctx.vertexAttributeLocation, 2, 0)
	ctx.uvAttribute = GlTools.NewAttribute(ctx.uvAttributeLocation, 2, 2)
}

func (ctx *Context) Render() {
	var time = float32(glfw.GetTime())
	ctx.deltaTime = time - ctx.lastTime
	ctx.fps = 1 / ctx.deltaTime
	ctx.lastTime = time

	w32.SetWindowPos(
		ctx.hwnd,
		w32.HWND_TOPMOST,
		0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_SHOWWINDOW,
	)

	var modelMatrix mgl32.Mat4

	width, height := ctx.window.GetSize()

	ortho := mgl32.Ortho(0, float32(width), float32(height), 0, -10000, 10000) //mgl32.Ortho2D(0, float32(width), float32(height), 0)

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)
	gl.UseProgram(ctx.mainProgram)

	ctx.lineRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.lines {
		ctx.lineRenderObject.UploadMesh(Mesh.Line(v.X1, v.Y1, v.X2, v.Y2))
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
		modelMatrix = mgl32.Translate3D(v.TransformX, v.TransformY, float32(v.ZIndex))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform1i(ctx.textureEnabledUniform, 0)
		ctx.lineRenderObject.Render()
	}
	ctx.lineRenderObject.End()

	ctx.outlineRectRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.outlineRects {
		modelMatrix = mgl32.Translate3D(v.X, v.Y, float32(v.ZIndex)).Mul4(mgl32.Translate3D(v.TransformX, v.TransformY, 0)).Mul4(mgl32.HomogRotate3DZ(v.Rotation)).Mul4(mgl32.Translate3D(-v.AnchorPointX*v.Width, -v.AnchorPointY*v.Height, 0)).Mul4(mgl32.Scale3D(v.Width, v.Height, 1))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureEnabledUniform, 0)
		ctx.outlineRectRenderObject.Render()
	}
	ctx.outlineRectRenderObject.End()

	ctx.rectRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()

	for _, v := range ctx.rects {
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
		modelMatrix = mgl32.Translate3D(v.X, v.Y, float32(v.ZIndex)).Mul4(mgl32.Translate3D(v.TransformX, v.TransformY, 0)).Mul4(mgl32.HomogRotate3DZ(v.Rotation)).Mul4(mgl32.Translate3D(-v.AnchorPointX*v.Width, -v.AnchorPointY*v.Height, 0)).Mul4(mgl32.Scale3D(v.Width, v.Height, 1))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureEnabledUniform, 0)

		ctx.rectRenderObject.Render()

	}
	ctx.rectRenderObject.End()

	ctx.circleRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.circles {
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
		modelMatrix = mgl32.Translate3D(v.X, v.Y, float32(v.ZIndex)).Mul4(mgl32.Translate3D(v.TransformX, v.TransformY, 0)).Mul4(mgl32.HomogRotate3DZ(v.Rotation)).Mul4(mgl32.Translate3D(-v.AnchorPointX*v.ScaleX, -v.AnchorPointY*v.ScaleY, 0)).Mul4(mgl32.Scale3D(v.ScaleX, v.ScaleY, 1))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureEnabledUniform, 0)
		ctx.circleRenderObject.Render()
	}
	ctx.circleRenderObject.End()
	ctx.polygonRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.polygons {
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
		ctx.polygonRenderObject.UploadMesh(Mesh.Polygon(
			v.X1, v.Y1,
			v.X2, v.Y2,
			v.X3, v.Y3,
		))

		modelMatrix = mgl32.Translate3D(v.TransformX, v.TransformY, float32(v.ZIndex)).Mul4(mgl32.HomogRotate3DZ(v.Rotation))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureEnabledUniform, 0)
		ctx.polygonRenderObject.Render()

	}
	ctx.polygonRenderObject.End()

	ctx.imageRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.images {
		ctx.imageRenderObject.ChangeTexture(v.Image)
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
		modelMatrix = mgl32.Translate3D(v.X, v.Y, float32(v.ZIndex)).Mul4(mgl32.Translate3D(v.TransformX, v.TransformY, 0)).Mul4(mgl32.HomogRotate3DZ(v.Rotation)).Mul4(mgl32.Translate3D(-v.AnchorPointX*v.Width, -v.AnchorPointY*v.Height, 0)).Mul4(mgl32.Scale3D(v.Width, v.Height, 1))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureUniform, 0)
		gl.Uniform1i(ctx.textureEnabledUniform, 1)

		ctx.imageRenderObject.Render()
	}
	ctx.imageRenderObject.End()

	ctx.textRenderObject.Begin()
	ctx.vertexAttribute.Use()
	ctx.uvAttribute.Use()
	for _, v := range ctx.texts {
		fontInfo, ok := ctx.loadedFonts[v.FontTexture]
		if !ok {
			continue
		}
		ctx.textRenderObject.UploadMesh(Mesh.Text(v.Text, fontInfo.Size.X, fontInfo.Size.Y, fontInfo.FontSize, v.Interval))
		ctx.textRenderObject.ChangeTexture(v.FontTexture)
		if v.Fill {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}

		var MaxLength int
		var MaxStrings int
		strs := strings.Split(v.Text, "\n")
		MaxStrings = len(strs)
		for _, s := range strs {
			if strLen := len(s); strLen > MaxLength {
				MaxLength = strLen
			}
		}

		modelMatrix = mgl32.Translate3D(v.X, v.Y, float32(v.ZIndex)).Mul4(mgl32.Translate3D(v.TransformX, v.TransformY, 0)).Mul4(mgl32.HomogRotate3DZ(v.Rotation)).Mul4(mgl32.Translate3D(-v.AnchorPointX*((v.Width+v.Size)*float32(MaxLength)*v.Interval), -v.AnchorPointY*((v.Height+v.Size)*float32(MaxStrings)), 0)).Mul4(mgl32.Scale3D(v.Width+v.Size, v.Height+v.Size, 1))
		gl.UniformMatrix4fv(ctx.cameraUniform, 1, false, &ortho[0])
		gl.UniformMatrix4fv(ctx.modelUniform, 1, false, &modelMatrix[0])
		gl.Uniform4fv(ctx.colorUniform, 1, &v.Color[0])
		gl.Uniform1i(ctx.textureUniform, 0)
		gl.Uniform1i(ctx.textureEnabledUniform, 1)

		ctx.textRenderObject.Render()
	}
	ctx.textRenderObject.End()

	gl.UseProgram(0)
	glfw.PollEvents()
	ctx.window.SwapBuffers()
}
func (ctx *Context) ClearAll() {
	ctx.lines = []Line{}
	ctx.rects = []Rect{}
	ctx.images = []Image{}
	ctx.polygons = []Polygon{}
	ctx.circles = []Circle{}
	ctx.outlineRects = []OutlineRect{}
	ctx.texts = []Text{}
}

func (ctx *Context) GetDeltaTime() float32 {
	return ctx.deltaTime
}
func (ctx *Context) GetFPS() float32 {
	return ctx.fps
}
func (ctx *Context) LoadFont(path string, FontSize float32) uint32 {
	atlas, err := ttf2atlas.FontToAtlas(path, FontSize)
	if err != nil {
		return 0
	}
	tex := GlTools.MakeTexture(true)
	w, h, _ := GlTools.UploadTextureFromImage(tex, atlas)
	ctx.loadedFonts[tex] = atlasFont{
		Size:     vec2{X: w, Y: h},
		FontSize: FontSize,
	}
	return tex
}
