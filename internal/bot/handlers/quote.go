package handlers

import (
	"log"

	"BastetTetlegram/internal/files"
	"BastetTetlegram/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QuoteHandler struct {
	bot *tgbotapi.BotAPI
}

func NewQuoteHandler(bot *tgbotapi.BotAPI) *QuoteHandler {
	return &QuoteHandler{
		bot: bot,
	}
}

func (h *QuoteHandler) Handle(update tgbotapi.Update) {
	log.Printf("Начата обработка команды /q для чата %d", update.Message.Chat.ID)

	phrases, err := files.ReadPhrasesFromFile("config/phrases.txt") // Или получите путь из конфига
	if err != nil {
		log.Printf("Ошибка при чтении файла фраз в команде /q: %v", err)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить цитату. Файл фраз недоступен."))
		return
	}

	if len(phrases) == 0 {
		log.Printf("Файл фраз пуст в команде /q для чата %d", update.Message.Chat.ID)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Файл с цитатами пуст."))
		return
	}

	randomPhrase := services.GetRandomPhrase(phrases)
	log.Printf("Выбрана случайная фраза: '%s'", randomPhrase)

	escapedPhrase := services.EscapeMarkdownV2(randomPhrase)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedPhrase)
	msg.ParseMode = "MarkdownV2"
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке цитаты в команде /q: %v", err)
	} else {
		log.Printf("Цитата успешно отправлена в чат %d", update.Message.Chat.ID)
	}
}
