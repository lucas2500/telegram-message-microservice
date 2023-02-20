package entities

type Message struct {
	BotToken       string `json:"BotToken"`
	ChatId         string `json:"ChatId"`
	Text           string `json:"message"`
	ParseMode      string `json:"ParseMode"`
	InlineKeyboard []struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data"`
	} `json:"Inline_Keyboard"`
	RetryAttempt int   `json:"RetryAttempt"`
	RetryOnError *bool `json:"RetryOnError"`
}

type Buttons []struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type ReplyMarkup struct {
	Buttons []Buttons `json:"inline_keyboard"`
}
