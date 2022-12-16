package handler

import (
	"fmt"
	"os"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"github.com/ted-vo/ilt-poker-club-bot/pkg/deck"
)

const (
	CURRENCY       = "💵"
	START          = "start"
	REGISTER       = "register"
	OPEN           = "open"
	CLOSE          = "close"
	ROLL           = "🎲 Roll"
	OPEN_DAILY     = "🎲 Open Daily Roll"
	OPEN_TOUR      = "🥇 Open Tournament Roll"
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
		tgbotapi.NewKeyboardButton(OPEN_DAILY),
		tgbotapi.NewKeyboardButton(OPEN_TOUR),
	),
)

var KeyboardPrivateButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PROFILE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(PERIODIC_TABLE),
		tgbotapi.NewKeyboardButton(HELP),
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
	case "draw_a_card":
		// lastDeck := deckMap[1]
		// if lastDeck == nil {
		// 	lastDeck := deck.NewDeck()
		// 	lastDeck.Shuffle()
		// 	msg.Text = "Open Deck and Shuffle"
		// 	deckMap[1] = lastDeck
		// } else {
		// 	card, err := lastDeck.Pop()
		// 	if err != nil {
		// 		log.Error(err.Error())
		// 	}
		// 	msg.Text = fmt.Sprintf("%s drawed: %s", handler.getCaller(update), card.ToString())
		// }
		// handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
	case OPEN_DAILY:
		handler.roll(DAILY_ROLL, update)
	case OPEN_TOUR:
		handler.roll(TOURNAMENT_ROLL, update)
	case PERIODIC_TABLE:
		handler.periodic_table(update)
	case PROFILE:
		handler.profile(update, &msg)
	case LEADERBOARD:
		// handler.leaderBoard(update, &msg)
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

	if update.Message.Chat.IsPrivate() || update.Message.Chat.IsChannel() {
		msg.Text = "This feature only for group!"
		handler.send(&msg)
		return
	}

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	switch rollType {
	case DAILY_ROLL:
		inlineKeyboard = rollDailyInlineKeyboard
	case TOURNAMENT_ROLL:
		inlineKeyboard = rollTournamentInlineKeyboard
	case "draw_a_card":
		inlineKeyboard = drawCardKeyboard
	}

	msg.ReplyMarkup = inlineKeyboard
	msg.Text = fmt.Sprintf("[ %s ] Ghi danh nào mấy con báo!", rollType)

	msgRes := handler.send(&msg)

	if rollType == "draw_a_card" {
		deck := deck.NewDeck()
		deck.Shuffle()
		deckMap[msgRes.MessageID] = deck
	}
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
