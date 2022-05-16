package entity

type Message struct {
	BotToken  string `json:"BotToken"`
	ChatId    string `json:"ChatId"`
	Text      string `json:"message"`
	ParseMode string `json:"ParseMode"`
}

type Response struct {
	Result string `json:"result"`
}
