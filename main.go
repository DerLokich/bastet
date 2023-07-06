package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1975934180:AAGxLBIi3FPEjjSz_V-YA0Z24JR_gPGY1JQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	substr := "сосед"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			messageText := update.Message.Text
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if strings.Contains(messageText, substr) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "YARRR!")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "NOPE")
				bot.Send(msg)
			}
		}
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		//UserName := update.Message.From.UserName
		cmdmsg := update.Message.MessageID
		//messageText := update.Message.Text
		//substr := "сосед"
		//if strings.Contains(messageText, substr) {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	msg.Text = "asdasdasd"
		//	bot.Send(msg)
		//} else {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	msg.Text = "NOPE"
		//	bot.Send(msg)
		//}

		if update.Message.Command() == "me" {
			//	msg.Text = "Надо бы удалить"
			//	bot.Send(msg)
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, cmdmsg)
			bot.Request(kill)
		}
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	if UserName == "DerLokich" {
		// Extract the command from the Message.
		//switch update.Message.Command() {
		//case "help":
		//	msg.Text = "I understand /sayhi and /status."
		//case "sayhi":
		//	msg.Text = "Hi :)"
		//case "getall":
		//	msg.Text = UserName
		//case "status":
		//	msg.Text = "I'm ok."
		//default:
		//	msg.Text = "Должен убить всех человеков..."
		//}
		//	}
		//	if _, err := bot.Send(msg); err != nil {
		//		log.Panic(err)
		//	}
	}
}
