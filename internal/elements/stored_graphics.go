package elements

const StoredGraphicsDefaultPath = "R:UNKNOWN.GRF"

type StoredGraphics struct {
	// The image data, in hexadecimal format. There is no default value.
	Data []byte

	// The total number of bytes in the image.
	// Because each pixel in the image uses 1 bit, this value should be the total number of pixels in the image, divided by 8 (since there are 8 bits per byte).
	// There is no default value.
	TotalBytes int

	// The number of bytes per pixel row in the image.
	// Because each pixel in the image uses 1 bit, this value should be the pixel width of the image, divided by 8 (since there are 8 bits per byte).
	// The default value is 1, which is almost always incorrect.
	RowBytes int
}
