package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("1975934180:AAGxLBIi3FPEjjSz_V-YA0Z24JR_gPGY1JQ"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
}
