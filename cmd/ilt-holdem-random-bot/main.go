package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-holdem-random-bot/pkg/handler"
)

func main() {
	log.Printf(os.Getenv("TELEGRAM_API_TOKEN"))
	bot, err := tgbotapi.NewBotAPI("5684691650:AAHjtjzAvHTaz8-PNMRUbwy5KT0I-vKkQ9U")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	u := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	u.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(u)

	handler := handler.New(bot)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				handler.Command(&update)
			}
		} else if update.CallbackQuery != nil {
			handler.InlineKeyboard(&update)
		}
	}
}
