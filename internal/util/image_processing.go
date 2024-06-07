package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func compressAndSaveImage(data []byte, quality int) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize the image to reduce its size
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3) // Resize to 800px width, preserving aspect ratio

	// Compress the image with specified quality
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}
