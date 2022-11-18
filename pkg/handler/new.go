package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func New(bot *tgbotapi.BotAPI) Handler {
	return &MessageHandler{
		bot: bot,
	}
}
