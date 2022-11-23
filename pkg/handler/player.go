package handler

import (
	"fmt"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"gopkg.in/Iwark/spreadsheet.v2"
	// "gopkg.in/Iwark/spreadsheet.v2"
)

const (
	SHEET_PLAYERS_TILE = "players"
	COL_ID             = 0
	COL_NAME           = 1
	COL_DEPOSIT        = 2
	COL_WITHDRAW       = 3
	COL_INCOME         = 4
	COL_RANK           = 5
)

type Player struct {
	Id       int64
	Name     string
	Desposit int64
	Withdraw int64
	Income   int64
}

func (handler *MessageHandler) registerPlayer(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerSheet, err := handler.SpreadsheetClub.Spreadsheet.SheetByTitle("players")
	if err != nil {
		log.Errorf("get players sheet error: %s", err.Error())
	}

	playerId := fmt.Sprintf("%v", update.Message.From.ID)
	playerName := handler.getCaller(update)

	isRegistered := isExisted(playerSheet.Rows, playerId)
	if isRegistered {
		msg.Text = "Bạn đã đăng kí rồi. Không thể đăng kí lại!"
		return
	}

	newRow := len(playerSheet.Rows)
	playerSheet.Update(newRow, COL_ID, playerId)
	playerSheet.Update(newRow, COL_NAME, fmt.Sprintf("%s", playerName))

	err = playerSheet.Synchronize()
	if err != nil {
		log.Errorf("sync player error: %s", err.Error())
	}

	text, _ := pkg.Parse("./config/register_success.html", struct{ ID string }{ID: playerId})
	msg.ParseMode = pkg.HTLM
	msg.Text = text
}

func isExisted(rows [][]spreadsheet.Cell, playerId string) bool {
	for _, row := range rows {
		if row[0].Value == playerId {
			return true
		}
	}

	return false
}
