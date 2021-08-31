package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1975934180:AAGxLBIi3FPEjjSz_V-YA0Z24JR_gPGY1JQ")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)
	for {
		select {
		case update := <-upd:
			// Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// Ответим пользователю его же сообщением
			reply := Text
			// Созадаем сообщение
			msg := tgbotapi.NewMessage(ChatID, reply)
			// и отправляем его
			bot.Send(msg)
		}

	}
}
