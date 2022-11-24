package handler

import (
	"fmt"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	case START:
		caller := handler.getCaller(update)
		msg.Text = fmt.Sprintf("Xin chào báo thủ %s đến với ILT Poker Club!", caller)
	case REGISTER:
		handler.registerPlayer(update, &msg)
	case OPEN:
		msg.Text = " 📜 Menu đã được thêm vào"
		msg.ReplyMarkup = KeyboardButton

		handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
	case CLOSE:
		msg.Text = " ❌  Loại bỏ Menu"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
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

func (handler *MessageHandler) getQuerier(update *tgbotapi.Update) string {
	querier := fmt.Sprintf("@%s", update.CallbackQuery.From.UserName)
	if len(querier) < 5 {
		querier = fmt.Sprintf("%s %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName)
	}

	return querier
}
