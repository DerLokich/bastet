package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1975934180:AAGxLBIi3FPEjjSz_V-YA0Z24JR_gPGY1JQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		UserName := update.Message.From.UserName
		cmdmsg := update.Message.MessageID
		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		if update.Message.Command() == "me" {
			//	msg.Text = "Надо бы удалить"
			//	bot.Send(msg)
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, cmdmsg)
			bot.Request(kill)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	if UserName == "DerLokich" {
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "getall":
			msg.Text = UserName
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "Должен убить всех человеков..."
		}
		//	}
		//	if _, err := bot.Send(msg); err != nil {
		//		log.Panic(err)
		//	}
	}
}
