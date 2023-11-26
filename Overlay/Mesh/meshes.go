package Mesh

import (
	"math"
	"strings"
)

/*
X,Y, U,V
*/
func Rect() []float32 {
	return []float32{
		0, 0, 0, 0,
		1, 0, 1, 0,
		1, 1, 1, 1,

		1, 1, 1, 1,
		0, 1, 0, 1,
		0, 0, 0, 0,
	}
}
func Line(X1, Y1, X2, Y2 float32) []float32 {
	return []float32{
		X1, Y1, 0, 0,

		X2, Y2, 1, 0,
	}
}

func Circle(numSegments int) []float32 {
	var vertices []float32
	radius := float32(0.5)

	for i := 0; i < numSegments; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numSegments)
		x := float32(math.Cos(angle)*float64(radius)) + radius
		y := float32(math.Sin(angle)*float64(radius)) + radius
		u := (x + radius) / (2 * radius)
		v := (y + radius) / (2 * radius)

		vertices = append(vertices, x, y, u, v)
	}
	return vertices
}
func Polygon(
	X1, Y1,
	X2, Y2,
	X3, Y3 float32,
) []float32 {
	return []float32{
		X1, Y1, 0, 0,
		X2, Y2, 0, 0,
		X3, Y3, 0, 0,
	}
}
func Text(str string, atlasWidth, atlasHeight int, fontSize, Interval float32) []float32 {
	var data []float32
	newLines := strings.Split(strings.ReplaceAll(str, "\t", " "), "\n")

	xSymbols := int32(float32(atlasWidth) / fontSize)
	normalFontSizeX := fontSize / float32(atlasWidth)
	normalFontSizeY := fontSize / float32(atlasHeight)

	for lineId, line := range newLines {
		y := float32(lineId)
		runes := []rune(line)
		for charId, char := range runes {
			texX := (float32(char%xSymbols) * fontSize) / float32(atlasWidth)
			texY := (float32(char/xSymbols) * fontSize) / float32(atlasHeight)
			x := float32(charId) * Interval
			data = append(data,
				x, y, texX, texY,
				1.0+x, y, normalFontSizeX+texX, texY,
				1.0+x, 1.0+y, normalFontSizeX+texX, normalFontSizeY+texY,

				1.0+x, 1.0+y, normalFontSizeX+texX, normalFontSizeY+texY,
				x, 1.0+y, texX, normalFontSizeY+texY,
				x, y, texX, texY,
			)
		}
	}

	return data
}
