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
	//Neib := ""
	substr := "сосед"
	Checker := 0
	LastMention := time.Now()
	//TimeDifference := time.Now()
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			messageText := update.Message.Text
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if strings.Contains(strings.ToLower(messageText), substr) {
				Checker++
				if LastMention != time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) {
					TimeDifference := time.Now().Sub(LastMention)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, TimeDifference.String())
					bot.Send(msg)
					LastMention = time.Now()
				}
			}
			//else {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "NOPE")
			//	bot.Send(msg)
			//}
		}
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		//UserName := update.Message.From.UserName
		cmdmsg := update.Message.MessageID

		if update.Message.Command() == "me" {
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, cmdmsg)
			bot.Request(kill)
		}

	}
}
