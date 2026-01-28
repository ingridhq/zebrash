package zebrash

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/internal/assets"
	"github.com/ingridhq/zebrash/internal/images"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

func TestImageRendering(t *testing.T) {
	font, err := truetype.Parse(assets.FontHelveticaBold)
	if err != nil {
		t.Fatal(err.Error())
	}

	var fontSize float64 = 18
	var scaleX float64 = float64(16) / float64(18)
	drawer := &drawer{
		matrix: gg.Identity(),
		face:   truetype.NewFace(font, &truetype.Options{Size: fontSize}),
		im:     newWhiteImage(813, 1626),
	}

	var x float64 = 734.5
	var y float64 = 80

	drawer.RotateAbout(gg.Radians(90), x, y)

	if scaleX != 1.0 {
		drawer.ScaleAbout(scaleX, 1, x, y)
	}

	drawer.drawString("SE-21000 MALMOE", x, y)

	fmt.Println("Pixel value at (743, 111)", drawer.im.At(743, 111))
}

func newWhiteImage(w, h int) draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for i := range img.Pix {
		img.Pix[i] = 255
	}

	return &imageWrap{img}
}

type drawer struct {
	matrix gg.Matrix
	face   font.Face
	im     draw.Image
}

func (drawer *drawer) drawString(s string, x, y float64) {
	d := &font.Drawer{
		Dst:  drawer.im,
		Src:  image.NewUniform(images.ColorBlack),
		Face: drawer.face,
		Dot:  fixp(x, y),
	}
	// based on Drawer.DrawString() in golang.org/x/image/font/font.go
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		sr := dr.Sub(dr.Min)
		transformer := draw.BiLinear
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := drawer.matrix.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		transformer.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})
		d.Dot.X += advance
		prevC = c
	}
}

func (drawer *drawer) Identity() {
	drawer.matrix = gg.Identity()
}

func (drawer *drawer) Translate(x, y float64) {
	drawer.matrix = drawer.matrix.Translate(x, y)
}

func (drawer *drawer) Scale(x, y float64) {
	drawer.matrix = drawer.matrix.Scale(x, y)
}

func (drawer *drawer) ScaleAbout(sx, sy, x, y float64) {
	drawer.Translate(x, y)
	drawer.Scale(sx, sy)
	drawer.Translate(-x, -y)
}

func (drawer *drawer) Rotate(angle float64) {
	drawer.matrix = drawer.matrix.Rotate(angle)
}

func (drawer *drawer) RotateAbout(angle, x, y float64) {
	drawer.Translate(x, y)
	drawer.Rotate(angle)
	drawer.Translate(-x, -y)
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fix(x), fix(y)}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

type imageWrap struct {
	img *image.RGBA
}

func (iw *imageWrap) ColorModel() color.Model {
	return iw.img.ColorModel()
}

func (iw *imageWrap) Bounds() image.Rectangle {
	return iw.img.Bounds()
}

func (iw *imageWrap) At(x, y int) color.Color {
	return iw.img.At(x, y)
}

func (iw *imageWrap) Set(x, y int, c color.Color) {
	if x == 743 && y == 111 {
		fmt.Println("Original pixel value at (743, 111)", c)
	}

	iw.img.Set(x, y, c)
}
