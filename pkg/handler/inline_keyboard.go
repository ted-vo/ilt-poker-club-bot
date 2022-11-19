package handler

import (
	"fmt"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var InlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	// tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData("Leader Board", "leaderboard"),
	// ),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 🎰 Roll", "roll"),
	),
)

func (hanlder *MessageHandler) InlineKeyboard(update *tgbotapi.Update) error {
	// Respond to the callback query, telling Telegram to show the user
	// a message with the data received.
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := hanlder.bot.Request(callback); err != nil {
		panic(err)
	}
	// And finally, send a message containing the data received.
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

	switch update.CallbackQuery.Data {
	case "roll":
		rolled := rand.Intn(12) + 1
		roller := fmt.Sprintf("@%s", update.CallbackQuery.From.UserName)
		if len(roller) == 0 {
			roller = fmt.Sprintf("%s %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName)
		}
		msg.Text = fmt.Sprintf("%s rolled: %d", roller, rolled)
	case "leaderboard":
		msg.Text = "Feature in development"
	}

	if _, err := hanlder.bot.Send(msg); err != nil {
		panic(err)
	}

	return nil
}
