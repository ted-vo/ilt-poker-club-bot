package handler

import (
	"fmt"
	"os"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
)

type Command interface {
	help(update *tgbotapi.Update, msg *tgbotapi.MessageConfig)
	menu(msg *tgbotapi.MessageConfig)
	periodic_table(mgs *tgbotapi.MessageConfig)
}

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
		handler.help(update, &msg)
	case "menu":
		handler.menu(&msg)
	case "periodictable":
		handler.periodic_table(update)
	default:
		msg.Text = "Tạm tời em không hiểu. Để em cập nhật thêm sau nhé!"
	}

	if len(msg.Text) != 0 {
		if _, err := handler.bot.Send(msg); err != nil {
			log.Error(err.Error())
		}
	}

	return nil
}

func (handler *MessageHandler) getCaller(update *tgbotapi.Update) string {
	caller := fmt.Sprintf("@%s", update.Message.From.UserName)
	if len(caller) == 0 {
		caller = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
	}

	return caller
}

func (handler *MessageHandler) help(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	caller := handler.getCaller(update)
	text, _ := pkg.Parse("./config/help.html",
		struct {
			Caller string
		}{
			Caller: caller,
		})
	msg.ParseMode = pkg.HTLM
	msg.Text = text
}

func (handler *MessageHandler) menu(msg *tgbotapi.MessageConfig) {
	msg.Text = " 🎲 Roll đi nào mấy con báo 🐆 "
	msg.ReplyMarkup = &InlineKeyboard
}

func (handler *MessageHandler) periodic_table(update *tgbotapi.Update) {
	bot := handler.bot
	f, err := os.Open("./config/periodic_table.jpg")
	if err != nil {
		log.Error(err.Error())
	}
	reader := tgbotapi.FileReader{Name: "periodic_table.jpg", Reader: f}
	msg := tgbotapi.NewPhoto(update.Message.Chat.ID, reader)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Caption = "Học đi nè con báo 🐆 "
	if _, err := bot.Send(msg); err != nil {
		log.Error(err.Error())
	}
}
