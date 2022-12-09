package resizer

import "github.com/KindCloud97/image_app/model"

type Resizer interface {
	Resize(img model.Image, percent uint8) (model.Image, error)
}
