// internal/bot/handlers/toast.go
package handlers

import (
	"log"

	"BastetTetlegram/internal/files"    // Импортируем пакет files
	"BastetTetlegram/internal/services" // Импортируем пакет services
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ToastHandler - структура обработчика команды /toast
type ToastHandler struct {
	bot *tgbotapi.BotAPI
	// Конфигурация может быть передана сюда, если нужно
	// config *config.Config
}

// NewToastHandler - конструктор для ToastHandler
func NewToastHandler(bot *tgbotapi.BotAPI /*, config *config.Config*/) *ToastHandler {
	return &ToastHandler{
		bot: bot,
		// config: config, // Если используется
	}
}

// Handle - метод, обрабатывающий команду /toast
func (h *ToastHandler) Handle(update tgbotapi.Update) {
	log.Printf("Начата обработка команды /toast для чата %d", update.Message.Chat.ID)

	// Вызов ПАКЕТНОЙ функции из internal/files
	toasts, err := files.ReadToastsFromFile("config/toasts.txt") // Или используйте путь из конфига
	if err != nil {
		log.Printf("Ошибка при чтении файла тостов в команде /toast: %v", err)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить тост. Файл тостов недоступен."))
		return
	}

	if len(toasts) == 0 {
		log.Printf("Файл тостов пуст в команде /toast для чата %d", update.Message.Chat.ID)
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Файл с тостами пуст."))
		return
	}

	// Вызов ПАКЕТНОЙ функции из internal/services
	randomToast := services.GetRandomToast(toasts)
	log.Printf("Выбран случайный тост: '%s'", randomToast)

	// Вызов ПАКЕТНОЙ функции из internal/services
	randomEmoji := services.GetRandomEmoji()
	log.Printf("Выбрано случайное эмодзи: '%s'", randomEmoji)

	// Вызов ПАКЕТНОЙ функции из internal/services
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
