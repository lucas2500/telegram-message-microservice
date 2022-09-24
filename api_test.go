package main

import (
	"log"
	"telegram-message-microservice/entities"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
)

func TestEnviaMensagem(t *testing.T) {

	client := req.C()
	Retry := false

	body := entities.Message{
		BotToken:     "",
		ChatId:       "",
		Text:         "",
		ParseMode:    "html",
		RetryOnError: &Retry,
	}

	for i := 0; i < 999; i++ {

		resp, err := client.R().
			SetBody(body).Post("http://localhost:3001/api/SendMessage")

		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 201, resp.StatusCode)
	}
}
