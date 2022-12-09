package queue

import (
	"github.com/KindCloud97/image_app/model"
)

type Queue interface {
	Send(img model.Image) error
	Receive() (model.Image, error)
}
