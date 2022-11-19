package handler

import (
	"fmt"
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
		caller := fmt.Sprintf("@%s", update.Message.From.UserName)
		if len(caller) == 0 {
			caller = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		}
		text, _ := pkg.Parse("./config/help.html",
			struct {
				Caller string
			}{
				Caller: caller,
			})
		msg.ParseMode = pkg.HTLM
		msg.Text = text
	case "menu":
		msg.Text = " 🎲 Roll đi nào mấy con báo 🐆 "
		msg.ReplyMarkup = InlineKeyboard
	default:
		msg.Text = "Tạm tời em không hiểu. Để em cập nhật thêm sau nhé!"
	}

	if _, err := handler.bot.Send(msg); err != nil {
		log.Panic(err)
	}
	return nil
}
