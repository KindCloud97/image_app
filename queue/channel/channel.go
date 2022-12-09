package channel

import (
	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/queue"
)

var _ queue.Queue = (*Channel)(nil)

type Channel struct {
	ch chan model.Image
}

func New() *Channel {
	return &Channel{
		ch: make(chan model.Image),
	}
}

func (c *Channel) Send(img model.Image) error {
	c.ch <- img
	return nil
}

func (c *Channel) Receive() (model.Image, error) {
	return <-c.ch, nil
}
