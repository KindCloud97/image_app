package handler

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net/http"

	"github.com/KindCloud97/image_app/model"
	"github.com/KindCloud97/image_app/queue"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type postResponse struct {
	ID string `json:"id"`
}

type errResponse struct {
	Message string `json:"message"`
}

func Get(c *gin.Context) {
	id := c.Param("ID")
	quality := c.DefaultQuery("quality", "100")

	c.File("tmp/" + id + ":" + quality)
}

func Post(c *gin.Context, q queue.Queue) {
	var body bytes.Buffer
	if c.Request.ContentLength > 0 {
		body.Grow(int(c.Request.ContentLength))
	}

	_, err := io.Copy(&body, c.Request.Body)
	if err != nil {
		log.Print(err)
		sendError(c, http.StatusInternalServerError)
		return
	}

	img := model.Image{
		ID:   uuid.NewString(),
		Data: body.Bytes(),
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(img); err != nil {
		log.Print(err)
		sendError(c, http.StatusInternalServerError)
		return
	}

	if err := q.Send(img); err != nil {
		log.Print(err)
		sendError(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, postResponse{ID: img.ID})

}

func sendError(c *gin.Context, code int) {
	c.JSON(code, errResponse{Message: http.StatusText(code)})
}
