package nfntresize

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/resizer"

	"github.com/nfnt/resize"
)

var _ resizer.Resizer = (*NfntResizer)(nil)

type NfntResizer struct {
}

func New() *NfntResizer {
	return &NfntResizer{}
}

func (r *NfntResizer) Resize(img model.Image, percent uint8) (model.Image, error) {
	file := bytes.NewBuffer(img.Data)

	// decode jpeg into image.Image
	data, _, err := image.Decode(file)
	if err != nil {
		return model.Image{}, err
	}

	width := uint(float64(data.Bounds().Dx()) * float64(percent) / 100.0)
	height := uint(float64(data.Bounds().Dy()) * float64(percent) / 100.0)

	minImg := resize.Resize(width, height, data, resize.Lanczos3)

	// write new image to file
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, minImg, nil)
	if err != nil {
		return model.Image{}, err
	}

	return model.Image{
		ID:   img.ID,
		Data: buf.Bytes(),
	}, nil
}
