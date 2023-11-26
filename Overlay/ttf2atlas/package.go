package ttf2atlas

import (
	"github.com/golang/freetype"
	"image"
	"os"
)

const defaultSymbolsPerWidth int = 32
const defaultSymbolsCount int = 2048

var symbolsPerWidth = defaultSymbolsPerWidth
var symbolsCount = defaultSymbolsCount

func SetSymbolsPerWidth(PerWidth int) { symbolsPerWidth = PerWidth }
func SetDefaultSymbolsPerWidth()      { symbolsPerWidth = defaultSymbolsPerWidth }
func SetSymbolsCount(Count int)       { symbolsCount = Count }
func SetDefaultSymbolsCount()         { symbolsCount = defaultSymbolsCount }

func FontToAtlas(FontPath string, FontSize float32) (*image.RGBA, error) {
	dpi := 72.0
	symbolSize := int(FontSize * float32(dpi) / 72)

	width := symbolsPerWidth * symbolSize
	height := (symbolsCount / symbolsPerWidth) * symbolSize

	data, err := os.ReadFile(FontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(data)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetFontSize(float64(FontSize))
	ctx.SetDPI(dpi)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.White)
	var y, x int
	for i := 0; i < symbolsCount; i++ {
		x = (i % symbolsPerWidth) * symbolSize
		y = (i / symbolsPerWidth) * symbolSize
		_, err = ctx.DrawString(string(rune(i)), freetype.Pt(x, y+symbolSize))
	}
	return img, nil
}
