// internal/bot/bot.go
package bot

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	handlers_pkg "BastetTetlegram/internal/bot/handlers" // Переименовали для избежания конфликта с builtin 'handlers'
	"BastetTetlegram/internal/config"
	"BastetTetlegram/internal/files"
	"BastetTetlegram/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tg     *tgbotapi.BotAPI
	cfg    *config.Config
	openai *services.OpenAIService
	// Добавьте другие сервисы
}

func New(cfg *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	openaiService := services.NewOpenAIService(cfg.OpenAI.Token)

	return &Bot{
		tg:     bot,
		cfg:    cfg,
		openai: openaiService,
	}, nil
}

func (b *Bot) Run() {
	// Загрузка LastMention
	LastMention, err := files.LoadLastMentionFromFile(b.cfg.Storage.LastMentionFile)
	if err != nil {
		if os.IsNotExist(err) {
			LastMention = time.Now()
			log.Printf("Файл с временем не найден, инициализация LastMention на текущее время: %v", LastMention)
		} else {
			log.Printf("Ошибка загрузки времени из файла, используется текущее время: %v", err)
			LastMention = time.Now()
		}
	} else {
		if LastMention.After(time.Now()) {
			log.Printf("Загруженное время в будущем, устанавливаем на текущее время.")
			LastMention = time.Now()
		}
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tg.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		messageText := update.Message.Text
		command := update.Message.Command()

		if command != "" {
			log.Printf("Получена команда: /%s от пользователя %d в чате %d", command, update.Message.From.ID, update.Message.Chat.ID)
		}

		// Регистрация обработчиков
		switch command {
		case "start":
			handler := handlers_pkg.NewStartHandler(b.tg, b.cfg)
			handler.Handle(update)
		case "help":
			handler := handlers_pkg.NewHelpHandler(b.tg, b.cfg)
			handler.Handle(update)
		case "q":
			handler := handlers_pkg.NewQuoteHandler(b.tg, b.cfg)
			handler.Handle(update)
		case "toast":
			handler := handlers_pkg.NewToastHandler(b.tg, b.cfg)
			handler.Handle(update)
		// case "gpt": ...
		// case "imagine": ...
		// case "me": ...
		// case "iddqd": ...
		default:
			// Логика для неизвестных команд или проверки "соседа"
		}

		if strings.Contains(strings.ToLower(messageText), "сосед") {
			// ... логика обновления LastMention и сохранения
			TimeDifference := time.Since(LastMention).Hours() / 24
			titles := []string{"день", "дня", "дней"}
			Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
			b.tg.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Neib))
			log.Println(TimeDifference)
			log.Printf("Предыдущее LastMention: %v", LastMention)
			LastMention = time.Now()
			log.Printf("Новое LastMention: %v", LastMention)

			err := files.SaveLastMentionToFile(b.cfg.Storage.LastMentionFile, LastMention)
			if err != nil {
				log.Printf("Ошибка сохранения времени в файл: %v", err)
			}
		}
	}
}

// declOfNum нужно будет переместить или импортировать
func declOfNum(number int, titles []string) string {
	if number < 0 {
		number *= -1
	}
	cases := []string{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}
