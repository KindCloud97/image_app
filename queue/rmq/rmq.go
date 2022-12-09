package rmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ queue.Queue = (*Rmq)(nil)

type Rmq struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}

func New(ch *amqp.Channel, msgs <-chan amqp.Delivery) *Rmq {
	return &Rmq{
		ch:   ch,
		msgs: msgs,
	}
}

func (r *Rmq) Send(img model.Image) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(img)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.ch.PublishWithContext(
		ctx,
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "",
			Body:        buf.Bytes(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rmq) Receive() (model.Image, error) {
	data := <-r.msgs
	buf := bytes.NewBuffer(data.Body)
	dec := gob.NewDecoder(buf)

	img := model.Image{}
	if err := dec.Decode(&img); err != nil {
		return model.Image{}, err
	}
	data.Ack(true)

	return img, nil
}
