package handler

import (
	"fmt"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"github.com/ted-vo/ilt-poker-club-bot/pkg/deck"
)

var rollMap = make(map[int]map[int64]*Roller)
var deckMap = make(map[int]*deck.Deck)

type Roller struct {
	Username     string
	Name         string
	RolledNumber int
	DrawedCard   *deck.Card
}

type RollType string

const (
	DAILY_ROLL      RollType = "🎲  Daily"
	TOURNAMENT_ROLL RollType = "🥇 Tournament"

	QUERY_DATA_DAILY_ROLL        = "daily_roll"
	QUERY_DATA_DAILY_ROLL_FINISH = "daily_roll_finish"
	QUERY_DATA_DAILY_ROLL_RESET  = "daily_roll_reset"

	QUERY_DATA_TOUR_ROLL        = "tour_roll"
	QUERY_DATA_TOUR_ROLL_FINISH = "tour_roll_finish"
	QUERY_DATA_TOUR_ROLL_RESET  = "tour_roll_reset"

	QUERY_DATA_DEPOSIT  = "deposit"
	QUERY_DATA_WITHDRAW = "withdraw"
)

var rollDailyInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ROLL, QUERY_DATA_DAILY_ROLL),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(FINISH, QUERY_DATA_DAILY_ROLL_FINISH),
		tgbotapi.NewInlineKeyboardButtonData(RESET, QUERY_DATA_DAILY_ROLL_RESET),
	),
)

var rollTournamentInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ROLL, QUERY_DATA_TOUR_ROLL),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(FINISH, QUERY_DATA_TOUR_ROLL_FINISH),
		tgbotapi.NewInlineKeyboardButtonData(RESET, QUERY_DATA_TOUR_ROLL_RESET),
	),
)

var drawCardKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🃏 Draw a card", "draw_a_card"),
	),
)

var walletInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(DEPOSIT, QUERY_DATA_DEPOSIT),
		tgbotapi.NewInlineKeyboardButtonData(WITHDRAW, QUERY_DATA_WITHDRAW),
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
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	switch update.CallbackQuery.Data {
	case "draw_a_card":
		hanlder.draw_card_query(update, &msg)
	// Daily
	case QUERY_DATA_DAILY_ROLL:
		hanlder.roll_query(DAILY_ROLL, update, &msg)
	case QUERY_DATA_DAILY_ROLL_FINISH:
		hanlder.roll_query_finish(DAILY_ROLL, update, &msg)
	case QUERY_DATA_DAILY_ROLL_RESET:
		hanlder.roll_query_reset(DAILY_ROLL, update, &msg)

		// Tournament
	case QUERY_DATA_TOUR_ROLL:
		hanlder.roll_query(TOURNAMENT_ROLL, update, &msg)
	case QUERY_DATA_TOUR_ROLL_FINISH:
		hanlder.roll_query_finish(TOURNAMENT_ROLL, update, &msg)
	case QUERY_DATA_TOUR_ROLL_RESET:
		hanlder.roll_query_reset(TOURNAMENT_ROLL, update, &msg)

	case QUERY_DATA_DEPOSIT:
		hanlder.deposit(update, &msg)
	case QUERY_DATA_WITHDRAW:
		hanlder.withdraw(update, &msg)
	}

	if len(msg.Text) != 0 {
		if _, err := hanlder.bot.Send(msg); err != nil {
			panic(err)
		}
	}

	return nil
}

func (handler *MessageHandler) draw_card_query(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	rollerId := update.CallbackQuery.From.ID

	groupRollMap := rollMap[messageId]
	if groupRollMap == nil {
		handler.removeMessage(chatId, messageId)
		return
	}

	deck := deckMap[messageId]
	if deck == nil {
		handler.removeMessage(chatId, messageId)
		return
	}

	if roller := groupRollMap[rollerId]; roller != nil {
		msg.Text = fmt.Sprintf("%s roll rồi thì ngồi im đi con báo này!", roller.showName())
		return
	}

	card, err := deck.Pop()
	if err != nil {
		log.Error(err.Error())
	}
	groupRollMap[update.CallbackQuery.From.ID] = &Roller{
		Username:   fmt.Sprintf("@%s", update.CallbackQuery.From.UserName),
		Name:       fmt.Sprintf("%s %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName),
		DrawedCard: card,
	}

	text := fmt.Sprintf("Hãy rút cho mình 1 lá bài may mắn nào mấy con báo!\n\n")
	var index = 0
	for _, v := range groupRollMap {
		text += fmt.Sprintf("%d. %s\n", index+1, v.parseDrawedText())
		index++
	}

	editMessage := tgbotapi.NewEditMessageTextAndMarkup(
		chatId,
		messageId,
		text,
		drawCardKeyboard)
	handler.bot.Send(editMessage)

	return

}

func (handler *MessageHandler) roll_query(rollType RollType, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	rollerId := update.CallbackQuery.From.ID

	groupRollMap := rollMap[messageId]
	if groupRollMap == nil {
		handler.removeMessage(chatId, messageId)
		return nil
	}

	if roller := groupRollMap[rollerId]; roller != nil {
		msg.Text = fmt.Sprintf("%s roll rồi thì ngồi im đi con báo này!", roller.showName())
		return nil
	}

	groupRollMap[update.CallbackQuery.From.ID] = &Roller{
		Username:     fmt.Sprintf("@%s", update.CallbackQuery.From.UserName),
		Name:         fmt.Sprintf("%s %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName),
		RolledNumber: pkg.Rollem(),
	}

	text := fmt.Sprintf("[ %s ] Ghi danh nào mấy con báo!\n\n", rollType)
	var index = 0
	for _, v := range groupRollMap {
		text += fmt.Sprintf("%d. %s\n", index+1, v.parseRolledText())
		index++
	}

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	if rollType == DAILY_ROLL {
		inlineKeyboard = rollDailyInlineKeyboard
	} else {
		inlineKeyboard = rollTournamentInlineKeyboard
	}

	editMessage := tgbotapi.NewEditMessageTextAndMarkup(
		chatId,
		messageId,
		text,
		inlineKeyboard)
	handler.bot.Send(editMessage)

	return nil
}

func (handler *MessageHandler) roll_query_finish(rollType RollType, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID

	groupRollMap := rollMap[messageId]
	if groupRollMap == nil {
		handler.removeMessage(chatId, messageId)
		return nil
	}

	text := fmt.Sprintf("[ %s ] Những con báo sau chuẩn bị tinh thần quay lô nào!\n\n", rollType)
	var index = 0
	for _, v := range groupRollMap {
		text += fmt.Sprintf("%d. %s\n", index+1, v.parseRolledText())
		index++
	}

	editMessage := tgbotapi.NewEditMessageText(
		chatId,
		messageId,
		text,
	)
	handler.bot.Send(editMessage)

	return nil
}

func (handler *MessageHandler) roll_query_reset(rollType RollType, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID

	groupRollMap := rollMap[messageId]
	if groupRollMap == nil {
		handler.removeMessage(chatId, messageId)
		return nil
	}

	text := fmt.Sprintf("[ %s ] Ghi danh nào mấy con báo!\n\n", rollType)

	editMessage := tgbotapi.NewEditMessageTextAndMarkup(
		chatId,
		messageId,
		text,
		rollTournamentInlineKeyboard,
	)
	handler.bot.Send(editMessage)

	return nil
}

func (roller *Roller) showName() string {
	showName := roller.Username
	if len(showName) < 5 {
		showName = roller.Name
	}

	return showName
}

func (roller *Roller) parseRolledText() string {
	return fmt.Sprintf("%s rolled: %d", roller.showName(), roller.RolledNumber)
}

func (roller *Roller) parseDrawedText() string {
	return fmt.Sprintf("%s drawed: %s", roller.showName(), roller.DrawedCard.ToString())
}
