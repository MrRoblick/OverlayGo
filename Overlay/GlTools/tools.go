package GlTools

import (
	"DrawerGO/Overlay/Engine"
	"DrawerGO/Overlay/Type"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

type Attribute struct {
	attribLocation uint32
	offset, size   int32
}
type RenderObject struct {
	vao, vbo, textureId uint32
	dataLength          int32
	mode                string
}

func (obj *RenderObject) Begin() {
	BeginBuffers(obj.vao, obj.vbo, obj.textureId)

}
func (obj *RenderObject) Render() {
	switch obj.mode {
	case Type.Line:
		{
			gl.DrawArrays(gl.LINES, 0, obj.dataLength)
		}
	case Type.Rectangle:
		{
			gl.DrawArrays(gl.TRIANGLES, 0, obj.dataLength)
		}
	case Type.Text:
		{
			gl.DrawArrays(gl.TRIANGLES, 0, obj.dataLength)
		}
	case Type.Image:
		{
			gl.DrawArrays(gl.TRIANGLES, 0, obj.dataLength)
		}
	case Type.Polygon:
		{
			gl.DrawArrays(gl.TRIANGLES, 0, obj.dataLength)
		}
	case Type.Circle:
		{
			gl.DrawArrays(gl.TRIANGLE_FAN, 0, obj.dataLength)
		}
	case Type.OutlineRect:
		{
			gl.DrawArrays(gl.LINE_LOOP, 0, obj.dataLength)
		}
	}
}

func (obj *RenderObject) End() {
	EndBuffers()
}
func (obj *RenderObject) UploadMesh(data []float32) {
	UploadToArrayBuffer(obj.vbo, data)
	obj.dataLength = int32(len(data))
	data = nil
}

func (obj *RenderObject) Delete() {
	EndBuffers()
	DeleteBuffers(obj.vao, obj.vbo, obj.textureId)
}
func (attrib *Attribute) Use() {
	gl.EnableVertexAttribArray(attrib.attribLocation)
	gl.VertexAttribPointerWithOffset(attrib.attribLocation, attrib.size, gl.FLOAT, false, int32(Engine.VertexSize)*int32(Engine.FloatSize), uintptr(attrib.offset*int32(Engine.FloatSize)))
}
func (obj *RenderObject) ChangeTexture(TextureId uint32) {
	obj.textureId = TextureId
}

func NewRenderObject(Mode string) RenderObject {
	vao, vbo, textureId := MakeBuffers(false)
	return RenderObject{
		vao:        vao,
		vbo:        vbo,
		dataLength: 0,
		textureId:  textureId,
		mode:       Mode,
	}
}
func NewAttribute(AttributeLocation uint32, size, offset int32) Attribute {
	return Attribute{
		attribLocation: AttributeLocation,
		offset:         offset,
		size:           size,
	}
}

func NewProgram(Vertex, Fragment uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, Vertex)
	gl.AttachShader(program, Fragment)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("link program error: %v", log)
	}

	return program, nil
}

func CompileShader(SourceCode string, ShaderType uint32) (uint32, error) {
	shader := gl.CreateShader(ShaderType)

	csources, free := gl.Strs(SourceCode + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}
func DeleteBuffers(vao, vbo, textureId uint32) {
	if vao != 0 {
		gl.DeleteVertexArrays(1, &vao)
	}
	if vbo != 0 {
		gl.DeleteBuffers(1, &vbo)
	}
	if textureId != 0 {
		gl.DeleteTextures(1, &textureId)
	}
}
func MakeBuffers(textureEnabled bool) (uint32, uint32, uint32) {
	var vao, vbo, textureId uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	if textureEnabled {
		textureId = MakeTexture(true)
	}

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return vao, vbo, textureId
}

func MakeTexture(linear bool) uint32 {
	method1 := gl.NEAREST
	method2 := gl.NEAREST
	if linear {
		method1 = gl.LINEAR
		method2 = gl.LINEAR
	}

	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(method1))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(method2))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return textureId
}
func UploadTexture(textureId uint32, TexturePath string) (int, int, error) {
	file, err := os.Open(TexturePath)
	defer file.Close()
	if err != nil {
		return 0, 0, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}
	return UploadTextureFromImage(textureId, img)
}
func UploadTextureFromImage(textureId uint32, Img image.Image) (int, int, error) {
	rgba := image.NewRGBA(Img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), Img, image.Point{}, draw.Src)
	size := rgba.Rect.Size()

	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return size.X, size.Y, nil
}

func BeginBuffers(Vao, Vbo, TextureId uint32) {
	gl.BindVertexArray(Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, Vbo)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, TextureId)

}

func EndBuffers() {
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func UploadToArrayBuffer(Vbo uint32, data []float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*Engine.FloatSize, gl.Ptr(data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	data = nil
}
