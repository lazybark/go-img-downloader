package imgwork

import (
	"fmt"
	"os"
)

// SaveImage writes image to  the file
func SaveImage(imageBytes []byte, fileName string) error {
	flags := os.O_CREATE | os.O_TRUNC
	f, err := os.OpenFile(fileName, flags, 0666)
	if err != nil {
		return fmt.Errorf("[SaveImage] %w", err)
	}

	_, err = f.Write(imageBytes)
	if err != nil {
		return fmt.Errorf("[SaveImage] %w", err)
	}
	f.Close()

	return nil
}
