package imgwork

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/sunshineplan/imgconv"
)

// ConvertToPNG uses buffer to convert decoded image into PNG
func ConvertToPNG(decoded image.Image, buffer ImgWriter) (image.Image, error) {
	newImg := image.NewRGBA(decoded.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), decoded, decoded.Bounds().Min, draw.Over)
	decoded = newImg

	err := imgconv.Write(&buffer, decoded, &imgconv.FormatOption{Format: imgconv.PNG})
	if err != nil {
		return nil, fmt.Errorf("convert to jpeg: %w", err)
	}

	decoded, _, err = buffer.Decode()
	if err != nil {
		return nil, fmt.Errorf("[ERROR DECODING] %s", err)
	}

	return decoded, nil
}
