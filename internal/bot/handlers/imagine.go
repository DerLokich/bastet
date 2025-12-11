package handlers

import (
	"context"
	"log"
	"time"

	"BastetTetlegram/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type ImagineHandler struct {
	bot       *tgbotapi.BotAPI
	openaiSvc *services.OpenAIService
}

func NewImagineHandler(bot *tgbotapi.BotAPI, openaiSvc *services.OpenAIService) *ImagineHandler {
	return &ImagineHandler{
		bot:       bot,
		openaiSvc: openaiSvc,
	}
}

func (h *ImagineHandler) Handle(update tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	respUrl, err := h.openaiSvc.CreateImage(ctx, openai.ImageRequest{ // Используем сервис
		Prompt:         update.Message.CommandArguments(),
		Size:           openai.CreateImageSize512x512, // Или из конфига
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	})
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось создать изображение. Пожалуйста, попробуйте позже."))
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, respUrl.Data[0].URL)
	h.bot.Send(msg)
}
