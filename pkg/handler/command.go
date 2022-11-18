package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
)

func (handler *MessageHandler) Command(update *tgbotapi.Update) error {
	if !update.Message.IsCommand() { // ignore any non-command Messages
		return nil
	}

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the update.Message.
	switch update.Message.Command() {
	case "help":
		text, _ := pkg.Parse("./config/help.html", struct{}{})
		msg.ParseMode = pkg.HTLM
		msg.Text = text
	case "menu":
		msg.Text = "Menu"
		msg.ReplyMarkup = InlineKeyboard
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := handler.bot.Send(msg); err != nil {
		log.Panic(err)
	}
	return nil
}
