package handler

import (
	"fmt"
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
	OPEN_TOUR      = "🥇 Open Tournament"
	FINISH         = "🏁 Finish"
	RESET          = "❌ Reset"
	DEPOSIT        = "💸 Deposit"
	WITHDRAW       = "💰 Withdraw"
	PERIODIC_TABLE = "📖 Periodic Table"
	PROFILE        = "👤 Profile"
	LEADERBOARD    = "🏆 Leaderboard"
	HELP           = "❓ Help"
	FEEDBACK       = "💡 Feedback"
)

var KeyboardButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ROLL),
		tgbotapi.NewKeyboardButton(OPEN_TOUR),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PROFILE),
		tgbotapi.NewKeyboardButton(LEADERBOARD),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PERIODIC_TABLE),
		tgbotapi.NewKeyboardButton(HELP),
		// tgbotapi.NewKeyboardButton(FEEDBACK),
	),
)

type Keyboard interface {
	openKeyboard(msg *tgbotapi.MessageConfig)
	closeKeyboard(msg *tgbotapi.MessageConfig)
	roll(update *tgbotapi.Update)
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
		handler.roll(DAILY_ROLL, update)
	case OPEN_TOUR:
		handler.roll(TOURNAMENT_ROLL, update)
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

func (handler *MessageHandler) roll(rollType RollType, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if !update.Message.Chat.IsGroup() {
		msg.Text = "This feature only for group!"
		handler.send(&msg)
		return
	}

	msg.ReplyMarkup = rollTournamentInlineKeyboard
	msg.Text = fmt.Sprintf("[ %s ] Ghi danh nào mấy con báo!", rollType)

	msgRes := handler.send(&msg)

	// init map for this messageID open tour with keyboard markup
	rollMap[msgRes.MessageID] = make(map[int64]*Roller)

	handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
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
