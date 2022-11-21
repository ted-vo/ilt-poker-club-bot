package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ted-vo/ilt-poker-club-bot/pkg"
	"github.com/ted-vo/ilt-poker-club-bot/pkg/handler"
)

func main() {
	log.SetHandler(pkg.NewLogHandler())
	go startHTTPServer()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Error(err.Error())
	}

	bot.Debug = os.Getenv("DEBUG") == "true"

	log.Infof("Authorized on account %s", bot.Self.UserName)

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
			log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				handler.Command(&update)
			} else {
				handler.Keyboard(&update)
			}
		} else if update.CallbackQuery != nil {
			handler.InlineKeyboard(&update)
		}
	}
}

func startHTTPServer() {
	log.Info("starting server...")
	http.HandleFunc("/", httpHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Infof("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err.Error())
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n Welcome to ILT Poker Club", name)
}
