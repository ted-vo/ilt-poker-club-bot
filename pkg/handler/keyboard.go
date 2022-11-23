package handler

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
)

const (
	CURRENCY       = "💵"
	START          = "start"
	REGISTER       = "register"
	OPEN           = "open"
	CLOSE          = "close"
	ROLL           = "🎲 Roll"
	PERIODIC_TABLE = "📖 Periodic Table"
	PROFILE        = "👤 Profile"
	LEADERBOARD    = "🏆 Leaderboard"
	HELP           = "❓ Help"
	FEEDBACK       = "💡 Feedback"
)

var KeyboardButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ROLL),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PERIODIC_TABLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PROFILE),
		tgbotapi.NewKeyboardButton(LEADERBOARD),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(HELP),
		tgbotapi.NewKeyboardButton(FEEDBACK),
	),
)

type Keyboard interface {
	openKeyboard(msg *tgbotapi.MessageConfig)
	closeKeyboard(msg *tgbotapi.MessageConfig)
	roll(update *tgbotapi.Update, msg *tgbotapi.MessageConfig)
	periodic_table(update *tgbotapi.Update)
	profile(update *tgbotapi.Update, msg *tgbotapi.MessageConfig)
	help(update *tgbotapi.Update, msg *tgbotapi.MessageConfig)
}

func (handler *MessageHandler) removeMessage(chatId int64, messageId int) {
	if _, err := handler.bot.Request(tgbotapi.NewDeleteMessage(chatId, messageId)); err != nil {
		log.Errorf("delete message erorr: %s", err.Error())
	}
}

func (handler *MessageHandler) Keyboard(update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	log.Debugf("%s", update.Message.Text)
	switch update.Message.Text {
	case ROLL:
		handler.roll(update, &msg)
	case PERIODIC_TABLE:
		handler.periodic_table(update)
	case PROFILE:
		handler.profile(update, &msg)
	case HELP:
		handler.help(update, &msg)
	}

	if len(msg.Text) != 0 {
		if _, err := handler.bot.Send(msg); err != nil {
			log.Error(err.Error())
		}
	}

	return nil
}

func (handler *MessageHandler) roll(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	rolled := rand.Intn(12) + 1
	roller := fmt.Sprintf("@%s", update.Message.From.UserName)
	if len(roller) < 5 {
		roller = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
	}
	msg.Text = fmt.Sprintf("%s rolled: %d", roller, rolled)

	// Remove request after rolled
	handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
}

func (handler *MessageHandler) profile(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	text, _ := pkg.Parse("./config/profile.html",
		struct {
			Name     string
			Deposit  string
			Withdraw string
			Income   string
		}{
			Name:     handler.getCaller(update),
			Deposit:  fmt.Sprintf("%d %s", 0, CURRENCY),
			Withdraw: fmt.Sprintf("%d %s", 0, CURRENCY),
			Income:   fmt.Sprintf("%d %s", 0, CURRENCY),
		})
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = pkg.HTLM
	msg.Text = text
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
