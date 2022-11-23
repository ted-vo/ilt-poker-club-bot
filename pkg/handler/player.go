package handler

import (
	"fmt"
	"strconv"

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
	Id       string
	Name     string
	Desposit int64
	Withdraw int64
	Income   int64
}

func getPlayer(rows [][]spreadsheet.Cell, playerId string) *Player {
	for _, row := range rows {
		if row[0].Value == playerId {
			desposit, _ := strconv.ParseInt(row[COL_DEPOSIT].Value, 10, 64)
			withdraw, _ := strconv.ParseInt(row[COL_DEPOSIT].Value, 10, 64)
			income, _ := strconv.ParseInt(row[COL_DEPOSIT].Value, 10, 64)
			return &Player{
				Id:       row[COL_ID].Value,
				Name:     row[COL_NAME].Value,
				Desposit: desposit,
				Withdraw: withdraw,
				Income:   income,
			}
		}
	}

	return nil
}

func (handler *MessageHandler) getPlayerSheet() *spreadsheet.Sheet {
	playerSheet, err := handler.SpreadsheetClub.Spreadsheet.SheetByTitle(SHEET_PLAYERS_TILE)
	if err != nil {
		log.Errorf("get players sheet error: %s", err.Error())
	}

	return playerSheet
}

func (handler *MessageHandler) sheetSync(sheet *spreadsheet.Sheet) {
	err := sheet.Synchronize()
	if err != nil {
		log.Errorf("sync error: %s", err.Error())
	}
}

func (handler *MessageHandler) registerPlayer(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerSheet := handler.getPlayerSheet()
	playerId := fmt.Sprintf("%v", update.Message.From.ID)
	playerName := handler.getCaller(update)

	msg.ReplyToMessageID = update.Message.MessageID

	if player := getPlayer(playerSheet.Rows, playerId); player != nil {
		msg.Text = "Bạn đã đăng kí rồi. Không thể đăng kí lại!"
		return
	}

	newRow := len(playerSheet.Rows)
	playerSheet.Update(newRow, COL_ID, playerId)
	playerSheet.Update(newRow, COL_NAME, fmt.Sprintf("%s", playerName))
	playerSheet.Update(newRow, COL_DEPOSIT, "0")
	playerSheet.Update(newRow, COL_WITHDRAW, "0")
	playerSheet.Update(newRow, COL_INCOME, "0")

	handler.sheetSync(playerSheet)

	text, err := pkg.Parse("./config/register_success.html", struct{ ID string }{ID: playerId})
	if err != nil {
		log.Error(err.Error())
	}
	msg.ParseMode = pkg.HTLM
	msg.Text = text
}

func (handler *MessageHandler) profile(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerSheet := handler.getPlayerSheet()
	playerId := fmt.Sprintf("%v", update.Message.From.ID)
	player := getPlayer(playerSheet.Rows, playerId)
	if player == nil {
		msg.Text = "Vui lòng đăng kí thông tin báo thủ!"
		return
	}

	text, _ := pkg.Parse("./config/profile.html",
		struct {
			Name     string
			Deposit  string
			Withdraw string
			Income   string
		}{
			Name:     handler.getCaller(update),
			Deposit:  fmt.Sprintf("%d %s", player.Desposit, CURRENCY),
			Withdraw: fmt.Sprintf("%d %s", player.Withdraw, CURRENCY),
			Income:   fmt.Sprintf("%d %s", player.Income, CURRENCY),
		})
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = pkg.HTLM
	msg.Text = text
}
