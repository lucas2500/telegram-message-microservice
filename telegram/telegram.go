package telegram

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"telegram-message-microservice/entities"
	"telegram-message-microservice/queue"
	"telegram-message-microservice/util"
)

func SendMessageToTelegram(body []byte) {

	var (
		message        entities.Message
		buttons        []entities.Buttons
		ReplyMarkup    entities.ReplyMarkup
		InlineKeyBoard []byte
		retry          bool = true
	)

	err := json.Unmarshal(body, &message)

	if err != nil {
		util.FailOnError(err, "Erro no parsing do json")
	}

	// Verifica se botões foram enviados no body
	if len(message.InlineKeyboard) > 0 {

		for i := range message.InlineKeyboard {
			button := entities.Buttons{{Text: message.InlineKeyboard[i].Text, CallbackData: message.InlineKeyboard[i].CallbackData}}
			buttons = append(buttons, button)
		}

		ReplyMarkup = entities.ReplyMarkup{Buttons: buttons}

		InlineKeyBoard, err = json.Marshal(&ReplyMarkup)

		if err != nil {
			util.FailOnError(err, "Erro ao serializar a mensagem")
		}
	}

	success := RequestTelegramAPI(message.BotToken, message.ChatId, message.Text, message.ParseMode, string(InlineKeyBoard))

	if message.RetryOnError != nil {
		retry = *message.RetryOnError
	}

	// Case o envio da mensagem falhe e a mensagem esteja parametrizada para reenvio, novas tentativas serão feitas
	if !success && retry {

		log.Println("Houve uma falha na tentativa", message.RetryAttempt, "de envio da mensagem!!")

		message.RetryAttempt++

		if message.RetryAttempt <= 500 {

			log.Println("Criando tentativa", message.RetryAttempt, "de envio da mensagem!!")
			CreateNewMessage(message)
			return
		}

		log.Println("Número máximo de tentativas excedido!! A mensagem será descartada")
	}
}

func RequestTelegramAPI(BotToken string, ChatId string, Text string, ParseMode string, InlineKeyBoard string) bool {

	req, err := http.Get(os.Getenv("TELEGRAM_BASE_URL") +
		BotToken + "/" + os.Getenv("TELEGRAM_ROUTE") + "?chat_id=" +
		ChatId + "&text=" + Text + "&parse_mode=" + ParseMode +
		"&reply_markup=" + InlineKeyBoard)

	if err != nil {
		log.Println("Erro interno ao enviar mensagem ao Telegram!!")
		return false
	}

	if req.StatusCode == 200 {
		log.Println("Mensagem enviada ao Telegram com sucesso!!")
		return true
	}

	log.Println("Erro ao enviar mensagem ao Telegram!!", "Status HTTP:", req.StatusCode)

	return false
}

func CreateNewMessage(message entities.Message) {

	j, err := json.Marshal(message)

	if err != nil {
		util.FailOnError(err, "Erro ao serializar a mensagem")
	}

	DLX := map[string]interface{}{
		"x-message-ttl":             10000,
		"x-dead-letter-exchange":    os.Getenv("RABBITMQ_MESSAGE_EXCHANGE"),
		"x-dead-letter-routing-key": os.Getenv("RABBITMQ_MESSAGE_QUEUE"),
	}

	QueuePros := queue.QueueProperties{
		Exchange:   os.Getenv("RABBITMQ_DELAY_MESSAGE_EXCHANGE"),
		RoutingKey: os.Getenv("RABBITMQ_DELAY_MESSAGE_ROUTING_KEY"),
		Queue:      os.Getenv("RABBITMQ_DELAY_MESSAGE_QUEUE"),
		DLX:        DLX,
		Body:       j,
	}

	if !QueuePros.QueueMessage() {
		log.Println("Houve um erro na criação da tentativa", message.RetryAttempt, "de envio da mensagem!!")
		return
	}

	log.Println("Nova tentativa de envio criada com sucesso!!")
}
