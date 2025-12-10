package images

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

func EncodeMonochrome(w io.Writer, img image.Image) error {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("img is not an RGBA image")
	}

	return png.Encode(w, sauvolaThresholding(rgba, 15, 0.5, 128))
}

// window = odd size of local neighborhood (e.g. 15 or 25)
// k = typically 0.2 to 0.5
// R = dynamic range of std-dev, usually 128
func sauvolaThresholding(src *image.RGBA, window int, k, R float64) *image.Gray {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	intImg := NewIntegralImage(w+1, h+1)
	intSqImg := NewIntegralImage(w+1, h+1)

	for y := range h {
		rowSum := 0.0
		rowSqSum := 0.0
		for x := range w {
			p := float64(src.RGBAAt(bounds.Min.X+x, bounds.Min.Y+y).R)
			rowSum += p
			rowSqSum += p * p

			intImg.Set(x+1, y+1, intImg.Get(x+1, y)+rowSum)
			intSqImg.Set(x+1, y+1, intSqImg.Get(x+1, y)+rowSqSum)
		}
	}

	dst := image.NewGray(bounds)

	r := window / 2

	for y := range h {
		for x := range w {
			x1 := max(0, x-r)
			y1 := max(0, y-r)
			x2 := min(w-1, x+r) + 1
			y2 := min(h-1, y+r) + 1

			sum := intImg.GetDiff(x1, y1, x2, y2)
			sumSq := intSqImg.GetDiff(x1, y1, x2, y2)
			area := float64((x2 - x1) * (y2 - y1))
			mean := sum / area
			variance := max(sumSq/area-mean*mean, 0)
			threshold := mean * (1 + k*(math.Sqrt(variance)/R-1))

			p := float64(src.RGBAAt(bounds.Min.X+x, bounds.Min.Y+y).R)

			if p > threshold {
				dst.SetGray(bounds.Min.X+x, bounds.Min.Y+y, color.Gray{Y: 255})
			} else {
				dst.SetGray(bounds.Min.X+x, bounds.Min.Y+y, color.Gray{Y: 0})
			}
		}
	}

	return dst
}

type IntegralImage struct {
	values []float64
	w      int
	h      int
}

func NewIntegralImage(w, h int) *IntegralImage {
	return &IntegralImage{
		w:      w,
		h:      h,
		values: make([]float64, w*h),
	}
}

func (img *IntegralImage) Get(x, y int) float64 {
	return img.values[y*img.w+x]
}

func (img *IntegralImage) Set(x, y int, value float64) {
	img.values[y*img.w+x] = value
}

func (img *IntegralImage) GetDiff(x1, y1, x2, y2 int) float64 {
	return img.Get(x2, y2) - img.Get(x2, y1) - img.Get(x1, y2) + img.Get(x1, y1)
}
