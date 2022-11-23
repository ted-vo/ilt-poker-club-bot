package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg/spreadsheet"
)

func New(bot *tgbotapi.BotAPI, sheetClub *spreadsheet.SpreadsheetClub) Handler {
	return &MessageHandler{
		bot:             bot,
		SpreadsheetClub: sheetClub,
	}
}
