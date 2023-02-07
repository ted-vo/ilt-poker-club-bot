package handler

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	SHEET_HISTORY_TILE = "history"

	ACTION_DEPOSIT  = "DEPOSIT"
	ACTION_WITHDRAW = "WITHDRAW"

	HISTORY_COL_ID     = 0
	HISTORY_COL_ACTION = 1
	HISTORY_COL_VALUE  = 2
	HISTORY_COL_DATE   = 3
)

var queueRequestInput = make(map[string]*Trasaction)

type Trasaction struct {
	Id     string
	Name   string
	Action string
	Value  int64
	Date   string
}

func (handler *MessageHandler) getHistorySheet() *spreadsheet.Sheet {
	historySheet, err := handler.SpreadsheetClub.Spreadsheet.SheetByTitle(SHEET_HISTORY_TILE)
	if err != nil {
		log.Errorf("get history sheet error: %s", err.Error())
	}

	return historySheet
}

func (handler *MessageHandler) deposit(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerId := fmt.Sprintf("%v", update.CallbackQuery.From.ID)
	transaction := &Trasaction{
		Id:     playerId,
		Name:   handler.getQuerier(update),
		Action: ACTION_DEPOSIT,
		Value:  0,
		Date:   update.CallbackQuery.Message.Time().Format("2006-January-02"),
	}
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Amount",
		Selective:             false,
	}
	msg.Text = "Muốn nạp bao nhiêu con báo?"
	queueRequestInput[playerId] = transaction
}

func (handler *MessageHandler) withdraw(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerId := fmt.Sprintf("%v", update.CallbackQuery.From.ID)
	transaction := &Trasaction{
		Id:     playerId,
		Name:   handler.getQuerier(update),
		Action: ACTION_WITHDRAW,
		Value:  0,
		Date:   update.CallbackQuery.Message.Time().Format("2006-January-02"),
	}
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Amount",
		Selective:             false,
	}
	msg.Text = "Muốn rút bao nhiêu con báo?"
	queueRequestInput[playerId] = transaction
}

func (handler *MessageHandler) Transaction(update *tgbotapi.Update) error {
	playerId := fmt.Sprintf("%v", update.Message.From.ID)
	transaction := queueRequestInput[playerId]
	if transaction == nil {
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	value, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		msg.Text = "Vui lập nhập định dạng số."
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply:            true,
			InputFieldPlaceholder: "Amount",
			Selective:             false,
		}
		handler.send(&msg)

		return nil
	}

	transaction.Value = value
	handler.addTransaction(transaction, &msg)

	return nil
}

func (handler *MessageHandler) addTransaction(trasaction *Trasaction, msg *tgbotapi.MessageConfig) {
	historySheet := handler.getHistorySheet()
	newRow := len(historySheet.Rows)

	historySheet.Update(newRow, HISTORY_COL_ID, trasaction.Id)
	historySheet.Update(newRow, HISTORY_COL_ACTION, trasaction.Action)
	historySheet.Update(newRow, HISTORY_COL_VALUE, fmt.Sprintf("%d", trasaction.Value))
	historySheet.Update(newRow, HISTORY_COL_DATE, trasaction.Date)

	handler.sheetSync(historySheet)
	queueRequestInput[trasaction.Id] = nil

	text, err := pkg.Parse("./config/add_history_success.html", struct {
		Action string
	}{
		Action: trasaction.Action,
	})
	if err != nil {
		log.Error(err.Error())
	}
	msg.ParseMode = pkg.HTLM
	msg.Text = text
	msg.ReplyMarkup = KeyboardPrivateButton

	handler.send(msg)
}
