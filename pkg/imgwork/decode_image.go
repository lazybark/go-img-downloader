package imgwork

import (
	"bytes"
	"fmt"
	"image"
)

// ImgWriter is an implementation of io.Writer that can store, update & decode image data
type ImgWriter struct {
	Bytes []byte
}

// Write writes bytes to image data
func (w *ImgWriter) Write(b []byte) (int, error) {
	w.Bytes = append(w.Bytes, b...)
	return len(b), nil
}

// Decode decodes image and returns Image interface & string format
func (w *ImgWriter) Decode() (image.Image, string, error) {
	return DecodeImage(w.Bytes)
}

// DecodeImage image and returns Image interface & string format
func DecodeImage(imageBytes []byte) (image.Image, string, error) {
	reader := bytes.NewReader(imageBytes)
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("[DecodeImage] decode: %w", err)
	}
	return img, format, nil
}
