package common

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

func Resize(file []byte, width int) ([]byte, error) {
	var src image.Image
	var output bytes.Buffer
	var err error

	mimetype := http.DetectContentType(file)

	r := bytes.NewReader(file)

	switch mimetype {
	case "image/jpeg":
		src, err = jpeg.Decode(r)
	case "image/png":
		src, err = png.Decode(r)
	case "image/webp":
		src, err = webp.Decode(r)
	}

	if err != nil {
		return []byte{}, err
	}

	if src.Bounds().Dx() <= width {
		return file, nil
	}

	ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	height := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	err = jpeg.Encode(&output, dst, nil)
	if err != nil {
		return []byte{}, err
	}

	return output.Bytes(), err
}
