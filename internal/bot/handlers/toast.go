// internal/bot/handlers/toast.go
package handlers

import (
	"log"

	"BastetTetlegram/internal/files"
	"BastetTetlegram/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ToastHandler struct {
	bot              *tgbotapi.BotAPI
	config           interface{}                // Замените на реальный тип конфига, если нужно
	fileService      *files.FileService         // Или передавайте напрямую функции ReadToastsFromFile
	generatorService *services.GeneratorService // Или передавайте напрямую функции GetRandom...
}

func NewToastHandler(bot *tgbotapi.BotAPI, config interface{}) *ToastHandler {
	return &ToastHandler{
		bot:    bot,
		config: config, // Передайте реальный конфиг, если он нужен
	}
}

func (h *ToastHandler) Handle(update tgbotapi.Update) {
	log.Printf("Начата обработка команды /toast для чата %d", update.Message.Chat.ID)

	// Предположим, у FileService есть метод ReadToastsFromFile
	// toasts, err := h.fileService.ReadToastsFromFile(h.config.Files.ToastsFile)
	// Или вызываем напрямую:
	toasts, err := files.ReadToastsFromFile("config/toasts.txt") // Или получите путь из конфига
	if err != nil {
		log.Printf("Ошибка при чтении файла тостов в команде /toast: %v", err)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить тост. Файл тостов недоступен."))
		return
	}

	if len(toasts) == 0 {
		log.Printf("Файл тостов пуст в команде /toast для чата %д", update.Message.Chat.ID)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Файл с тостами пуст."))
		return
	}

	randomToast := services.GetRandomToast(toasts)
	log.Printf("Выбран случайный тост: '%s'", randomToast)

	randomEmoji := services.GetRandomEmoji()
	log.Printf("Выбрано случайное эмодзи: '%s'", randomEmoji)

	escapedToast := services.EscapeMarkdownV2(randomToast)
	finalMessage := randomEmoji + " " + escapedToast + " " + randomEmoji

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalMessage)
	msg.ParseMode = "MarkdownV2"
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке тоста в команде /toast: %v", err)
	} else {
		log.Printf("Тост с эмодзи успешно отправлен в чат %d", update.Message.Chat.ID)
	}
}
