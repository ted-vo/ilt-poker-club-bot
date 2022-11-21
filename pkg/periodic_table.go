package pkg

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var PeriodicTableKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("AA"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("AKo"),
	),
)
