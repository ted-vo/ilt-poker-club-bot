package handler

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/apex/log"
	"github.com/aquasecurity/table"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	SHEET_PLAYERS_TILE = "players"

	PLAYER_COL_ID       = 0
	PLAYER_COL_NAME     = 1
	PLAYER_COL_DEPOSIT  = 2
	PLAYER_COL_WITHDRAW = 3
	PLAYER_COL_INCOME   = 4
	PLAYER_COL_RANK     = 5
)

type Player struct {
	Id       string
	Name     string
	Deposit  int64
	Withdraw int64
	Income   int64
	Rank     uint64
}

func getPlayer(rows [][]spreadsheet.Cell, playerId string) *Player {
	for _, row := range rows {
		if row[0].Value == playerId {
			desposit, _ := strconv.ParseInt(row[PLAYER_COL_DEPOSIT].Value, 10, 64)
			withdraw, _ := strconv.ParseInt(row[PLAYER_COL_WITHDRAW].Value, 10, 64)
			income, _ := strconv.ParseInt(row[PLAYER_COL_INCOME].Value, 10, 64)
			rank, _ := strconv.ParseUint(row[PLAYER_COL_RANK].Value, 10, 64)
			return &Player{
				Id:       row[PLAYER_COL_ID].Value,
				Name:     row[PLAYER_COL_NAME].Value,
				Deposit:  desposit,
				Withdraw: withdraw,
				Income:   income,
				Rank:     rank,
			}
		}
	}

	return nil
}

func getPlayers(rows [][]spreadsheet.Cell) []*Player {
	var players = make([]*Player, 0)
	for _, row := range rows {
		if row[PLAYER_COL_ID].Value == "id" {
			continue
		}

		desposit, _ := strconv.ParseInt(row[PLAYER_COL_DEPOSIT].Value, 10, 64)
		withdraw, _ := strconv.ParseInt(row[PLAYER_COL_WITHDRAW].Value, 10, 64)
		income, _ := strconv.ParseInt(row[PLAYER_COL_INCOME].Value, 10, 64)
		rank, _ := strconv.ParseUint(row[PLAYER_COL_RANK].Value, 10, 64)

		players = append(players, &Player{
			Id:       row[PLAYER_COL_ID].Value,
			Name:     row[PLAYER_COL_NAME].Value,
			Deposit:  desposit,
			Withdraw: withdraw,
			Income:   income,
			Rank:     rank,
		})
	}

	return players
}

func (handler *MessageHandler) getPlayerSheet() *spreadsheet.Sheet {
	handler.SpreadsheetClub.Reload()

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
	playerSheet.Update(newRow, PLAYER_COL_ID, playerId)
	playerSheet.Update(newRow, PLAYER_COL_NAME, fmt.Sprintf("%s", playerName))

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

	player.Name = handler.getCaller(update)
	text, _ := pkg.Parse("./config/profile.html", player)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = pkg.HTLM
	msg.Text = text

	// check private chat caller
	if update.Message.Chat.IsPrivate() {
		msg.ReplyMarkup = walletInlineKeyboard
	}
}

func (handler *MessageHandler) leaderBoard(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	playerSheet := handler.getPlayerSheet()
	players := getPlayers(playerSheet.Rows)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Income > players[j].Income
	})

	buf := new(bytes.Buffer)
	t := table.New(buf)
	t.SetHeaders("Rank", "Name", "Income")

	for _, v := range players {
		t.AddRow(fmt.Sprintf("%d", v.Rank), v.Name, fmt.Sprintf("%d", v.Income))
	}
	t.Render()
	data := buf.String()

	text, _ := pkg.Parse("./config/leaderboard.html", struct{ Data string }{Data: data})
	msg.ParseMode = pkg.HTLM
	msg.Text = text

	handler.removeMessage(update.Message.Chat.ID, update.Message.MessageID)
}
