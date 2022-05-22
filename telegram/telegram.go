package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"telegram-message-microservice/entity"
)

func SendMessageToTelegram(body []byte) {

	var message entity.Message
	var buttons []entity.Buttons

	err := json.Unmarshal(body, &message)

	if err != nil {
		log.Fatal(err, "Erro no parsing do json!!")
	}

	for i := range message.InlineKeyboard {
		button := entity.Buttons{{Text: message.InlineKeyboard[i].Text, CallbackData: message.InlineKeyboard[i].CallbackData}}
		buttons = append(buttons, button)
	}

	ReplyMarkup := entity.ReplyMarkup{Buttons: buttons}

	InlineKeyBoard, err := json.Marshal(&ReplyMarkup)

	if err != nil {
		log.Fatal(err, "Erro ao marshalar o json!!")
	}

	req, err := http.Get(os.Getenv("TELEGRAM_BASE_URL") + message.BotToken + "/" + os.Getenv("TELEGRAM_ROUTE") + "?chat_id=" + message.ChatId + "&text=" + message.Text + "&parse_mode=" + message.ParseMode + "&reply_markup=" + string(InlineKeyBoard))

	if err != nil {
		log.Fatal(err, "Erro ao enviar mensagem ao Telegram!!")
	}

	if req.StatusCode == 200 {
		fmt.Println("Mensagem enviada ao Telegram com sucesso!!")
	} else {
		fmt.Println("Erro ao enviar mensagem ao Telegram!!")
		fmt.Println(req.StatusCode)
	}
}
