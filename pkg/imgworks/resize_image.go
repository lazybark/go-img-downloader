package imgworks

import (
	"image"

	"github.com/disintegration/imaging"
)

// ResizeImage returns resized version of image
func ResizeImage(img image.Image, width int, height int, filter imaging.ResampleFilter) image.Image {
	return imaging.Resize(img, width, height, filter)
}
