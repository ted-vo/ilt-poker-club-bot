package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg/spreadsheet"
)

type Handler interface {
	Command(update *tgbotapi.Update) error
	Keyboard(update *tgbotapi.Update) error
	InlineKeyboard(update *tgbotapi.Update) error
}

type MessageHandler struct {
	bot             *tgbotapi.BotAPI
	SpreadsheetClub *spreadsheet.SpreadsheetClub
}
