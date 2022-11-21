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
	OPEN           = "!open"
	CLOSE          = "!close"
	ROLL           = "🎲 Roll"
	PREIODIC_TABLE = "📖 Priodic Table"
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
		tgbotapi.NewKeyboardButton(PREIODIC_TABLE),
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
}

func (handler *MessageHandler) removeMessage(chatId int64, messageId int) {
	handler.bot.Send(tgbotapi.NewDeleteMessage(chatId, messageId))
}

func (handler *MessageHandler) Keyboard(update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	log.Debugf("%s", update.Message.Text)
	switch update.Message.Text {
	case OPEN:
		msg.Text = " 📜 Menu đã được thêm vào"
		msg.ReplyMarkup = KeyboardButton
	case CLOSE:
		msg.Text = " ❌  Loại bỏ Menu"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case ROLL:
		handler.roll(update, &msg)
	case PREIODIC_TABLE:
		handler.periodic_table(update)
	case PROFILE:
		handler.profile(update, &msg)
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
