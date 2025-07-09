package elements

type GraphicFieldFormat int

const (
	GraphicFieldFormatHex GraphicFieldFormat = 1
	GraphicFieldFormatRaw GraphicFieldFormat = 2
	GraphicFieldFormatAR  GraphicFieldFormat = 3
)

type GraphicField struct {
	ReversePrint

	Position LabelPosition

	// The format of the image data contained in the fifth parameter.
	// Valid values are A (hexadecimal format), B (raw binary format), and C (AR compressed).
	// There is no default value.
	Format GraphicFieldFormat

	// The total number of data bytes in the fifth parameter.
	// The value of this parameter is always the same as totalBytes,
	// except in the case of format C (AR compressed) which is very rarely used.
	DataBytes int

	// The total number of bytes in the image.
	// Because each pixel in the image uses 1 bit, this value should be the total number of pixels in the image,
	// divided by 8 (since there are 8 bits per byte).
	// There is no default value.
	TotalBytes int

	// The number of bytes per pixel row in the image.
	// Because each pixel in the image uses 1 bit, this value should be the pixel width of the image,
	// divided by 8 (since there are 8 bits per byte).
	// There is no default value.
	RowBytes int

	// The image data, in the format specified in the first parameter.
	// There is no default value
	Data []byte

	// The horizontal magnification to apply to the image. Any number between 1 and 10 may be used. The default value is 1.
	MagnificationX int

	// The vertical magnification to apply to the image. Any number between 1 and 10 may be used. The default value is 1.
	MagnificationY int
}
