package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/KindCloud97/image_app/handler"
	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/queue"
	"github.com/KindCloud97/image_app/queue/rmq"
	"github.com/KindCloud97/image_app/resizer"
	"github.com/KindCloud97/image_app/resizer/nfntresize"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func worker(q queue.Queue, resizer resizer.Resizer) {
	for {
		img, err := q.Receive()
		if err != nil {
			log.Print(err)
			continue
		}

		err = ResizeAndSave(img, resizer)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func ResizeAndSave(img model.Image, resizer resizer.Resizer) error {
	for percent := 100; percent >= 25; percent -= 25 {
		img, err := resizer.Resize(img, uint8(percent))
		if err != nil {
			return err
		}

		filename := "tmp/" + img.ID + ":" + strconv.Itoa(percent)
		err = os.WriteFile(filename, img.Data, 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	r := gin.Default()

	q, cancel, err := rabbitMQQueue()
	if err != nil {
		log.Print(err)
	}
	defer cancel()

	resizer := nfntresize.New()
	for w := 1; w <= 3; w++ {
		go worker(q, resizer)
	}

	r.POST("/upload", func(c *gin.Context) {
		handler.Post(c, q)
	})

	r.GET("/:ID", func(c *gin.Context) {
		handler.Get(c)
	})

	r.Run()
}

func rabbitMQQueue() (queue.Queue, func(), error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"TestQueue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return rmq.New(ch, msgs),
		func() {
			conn.Close()
			ch.Close()
		},
		nil
}
