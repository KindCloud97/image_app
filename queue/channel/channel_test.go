package channel_test

import (
	"testing"

	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/queue/channel"
)

func TestChannel_Send_Receive(t *testing.T) {
	img := model.Image{
		ID:   "123",
		Data: []byte{123},
	}

	ch := channel.New()

	go func() {
		err := ch.Send(img)
		if err != nil {
			t.Fatal(err)
		}
	}()

	recievedImage, err := ch.Receive()
	if err != nil {
		t.Fatal(err)
	}
	if img.ID != recievedImage.ID {
		//TODO: Compare bytes
		t.Fatal("Not Equal")
	}
}
