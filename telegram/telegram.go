package telegram

import (
	"encoding/json"
	"fmt"
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

	// Case o envio da mensagem falhe e a mensagem esteja parametrizada para reenvio, novas tentativas serão feitas
	if !success && message.RetryOnError {

		fmt.Println("Houve uma falha na tentativa", message.RetryAttempt, "de envio da mensagem!!")

		message.RetryAttempt++

		if message.RetryAttempt <= 500 {

			fmt.Println("Criando tentativa", message.RetryAttempt, "de envio da mensagem!!")
			CreateNewMessage(message)
			return
		}

		fmt.Println("Número máximo de tentativas excedido!! A mensagem será descartada")
	}
}

func RequestTelegramAPI(BotToken string, ChatId string, Text string, ParseMode string, InlineKeyBoard string) bool {

	req, err := http.Get(os.Getenv("TELEGRAM_BASE_URL") + BotToken + "/" + os.Getenv("TELEGRAM_ROUTE") + "?chat_id=" + ChatId + "&text=" + Text + "&parse_mode=" + ParseMode + "&reply_markup=" + InlineKeyBoard)

	if err != nil {
		fmt.Println("Erro interno ao enviar mensagem ao Telegram!!")
		return false
	}

	if req.StatusCode == 200 {
		fmt.Println("Mensagem enviada ao Telegram com sucesso!!")
		return true
	}

	fmt.Println("Erro ao enviar mensagem ao Telegram!!", "Status HTTP:", req.StatusCode)

	return false
}

func CreateNewMessage(message entities.Message) {

	j, err := json.Marshal(message)

	if err != nil {
		util.FailOnError(err, "Erro ao serializar a mensagem")
	}

	exchange := os.Getenv("RABBITMQ_EXCHANGE_NAME")
	RoutingKey := os.Getenv("RABBITMQ_DLQ_ROUTING_KEY")
	QueueName := os.Getenv("RABBITMQ_DLQ_NAME")

	if !queue.QueueMessage(j, exchange, RoutingKey, QueueName) {
		fmt.Println("Houve um erro na criação da tentativa", message.RetryAttempt, "de envio da mensagem!!")
		return
	}

	fmt.Println("Nova tentativa de envio criada com sucesso!!")
}
