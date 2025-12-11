package handlers

import (
	"context"
	"log"
	"time"

	"BastetTetlegram/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type GPTHandler struct {
	bot       *tgbotapi.BotAPI
	openaiSvc *services.OpenAIService
	// req       openai.ChatCompletionRequest // Это состояние на чат/сессию, см. улучшения
}

func NewGPTHandler(bot *tgbotapi.BotAPI, openaiSvc *services.OpenAIService) *GPTHandler {
	return &GPTHandler{
		bot:       bot,
		openaiSvc: openaiSvc,
		// req: openai.ChatCompletionRequest{...} // Инициализация в конструкторе не подходит для глобального состояния
	}
}

func (h *GPTHandler) Handle(update tgbotapi.Update) {
	// Обратите внимание: текущая реализация использует глобальное состояние req.
	// Это приведет к проблемам при нескольких пользователях.
	// Для корректной работы нужен менеджер сессий (см. улучшения).
	// Пока оставим как есть, но с комментарием.

	// Временный фикс: создаем новый запрос для каждого вызова
	req := openai.ChatCompletionRequest{
		Temperature: 0.7,
		Model:       openai.GPT4o, // Или получите из конфига
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Temporary message for initialization", // Или из конфига
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: update.Message.CommandArguments(),
	})

	resp, err := h.openaiSvc.CreateChatCompletion(ctx, req) // Используем сервис
	if err != nil {
		apiErr, ok := err.(*openai.APIError)
		if ok && apiErr.HTTPStatusCode == 400 {
			cancel()
			// Повторный вызов с другой моделью
			req.Model = openai.GPT4oMini // Или получите из конфига
			req.Messages = []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Clear message",
				},
			}
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			})

			resp, err = h.openaiSvc.CreateChatCompletion(ctx, req) // Нужно создать новый контекст или отменить и создать
			if err != nil {
				// Обработка ошибки после переключения модели
				h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, что-то пошло не так. Пожалуйста, попробуйте снова."))
				log.Printf("Ошибка при вызове CreateChatCompletion после 400: %v\n", err)
				h.bot.Send(tgbotapi.NewMessage(435809098, err.Error())) // Или используйте ID из конфига
				return
			}
		} else {
			// Обработка других ошибок
			h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я устала запоминать, обнуляюсь"))
			log.Printf("Ошибка при вызове CreateChatCompletion: %v\n", err)
			h.bot.Send(tgbotapi.NewMessage(435809098, err.Error())) // Или используйте ID из конфига
			return
		}
	}

	// Отправка ответа
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Message.Content)
	h.bot.Send(msg)

	// Обратите внимание: добавление ответа в историю также ломается из-за отсутствия сессии.
	// В реальной реализации, req.Messages нужно хранить отдельно для каждого чата.
}
