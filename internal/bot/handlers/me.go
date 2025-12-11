package handlers

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MeHandler struct {
	bot *tgbotapi.BotAPI
}

func NewMeHandler(bot *tgbotapi.BotAPI) *MeHandler {
	return &MeHandler{
		bot: bot,
	}
}

func (h *MeHandler) Handle(update tgbotapi.Update) {
	time.Sleep(1 * time.Second) // Это плохая практика в обработчике, но в оригинале так
	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	}
	_, err := h.bot.Request(deleteMsg)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	}
}
