package main

import (
	"BastetTetlegram/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	substr := "сосед"
	Checker := 0
	LastMention := time.Now()
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
					TimeDifference := time.Since(LastMention).Hours() / 24
					//Days := int(TimeDifference)
					Neib := strconv.Itoa(int(TimeDifference)) + " дней без соседей"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, Neib)
					bot.Send(msg)
					log.Println(TimeDifference)
					log.Printf(LastMention.String())
					LastMention = time.Now()
					log.Printf(LastMention.String())
					log.Printf(Neib)
				}
			}
			//else {
			//	msg := tgbotapi.NewMessage(435809098, "NOPE")
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
