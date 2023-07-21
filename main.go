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

// func declOfNum принимает число и массив названий и возвращает строку с правильной формой названия
func declOfNum(number int, titles []string) string {
	// Если число отрицательное, приводим его к положительному
	if number < 0 {
		number *= -1
	}
	// Массив чисел для соответствия к каждому падежу
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	// Проверяем условия для определения падежа
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	// Возвращаем название соответствующего падежа
	return titles[currentCase]
}

func main() {
	titles := []string{"день", "дня", "дней"}
	LastMention := time.Now()
	//Создается экземпляр бота, используя токен, полученный из config.Token.
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	client := openai.NewClient(config.GPTtoken)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Temporary message for initialization",
			},
		},
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		if update.Message != nil { // If we got a message
			messageText := update.Message.Text
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Проверяет, содержит ли текст сообщения подстроку
			if strings.Contains(strings.ToLower(messageText), substr) {
				if LastMention != time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) {
					// Вычисляет разницу времени с момента последнего упоминания в днях
					TimeDifference := time.Since(LastMention).Hours() / 24
					// Создает сообщение с текстом, содержащим полученную разницу времени и отправляет его в чат
					Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Neib))
					bot.Send(tgbotapi.NewMessage(435809098, "Было: "+LastMention.String()))
					log.Println(TimeDifference)
					log.Printf(LastMention.String())
					LastMention = time.Now()
					log.Printf(LastMention.String())
					bot.Send(tgbotapi.NewMessage(435809098, "Стало: "+LastMention.String()))
				}
			}
		}
		// Данный фрагмент кода проверяет, является ли полученная команда от пользователя "me"
		cmdmsg := update.Message.MessageID
		if update.Message.Command() == "me" {
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, cmdmsg)
			bot.Request(kill)
		}

		// Этот фрагмент кода позволяет боту устанавливать определенные права доступа для указанного пользователя в чате при получении команды "iddqd".
		if update.Message.Command() == "iddqd" {
			// Создается переменная, которая используется для установки прав доступа для определенного пользователя в чате.
			memberConfig := tgbotapi.PromoteChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: -1001165249098,
					UserID: 435809098,
				},
				// Устанавливается значение true для разрешения выполнения соответствующих действий.
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
			// Выполняется запрос бота на изменение конфигурации пользователя
			bot.Request(memberConfig)
			// Отображается информация о memberConfig в журнале
			log.Println(memberConfig)
		}
		if update.Message.Command() == "gpt" {
			ctx := context.Background()
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			})
			// Создаем канал для отмены контекста (если понадобится)
			ctx, cancel := context.WithCancel(ctx)
			resp, err := client.CreateChatCompletion(ctx, req)
			if err != nil {
				// Проверяем, является ли ошибка ошибкой 400
				apiErr, ok := err.(*openai.APIError)
				if ok && apiErr.HTTPStatusCode == 400 {
					// Обрабатываем ошибку 400
					// Если ошибка является ошибкой 400, выполняются следующие действия: отменяется контекст, обновляется
					cancel()
					req = openai.ChatCompletionRequest{
						Model: openai.GPT3Dot5Turbo,
						Messages: []openai.ChatCompletionMessage{
							{
								Role:    openai.ChatMessageRoleSystem,
								Content: "Clear message",
							},
						},
					}
					errorDetails := apiErr.Error() // Получаем подробности ошибки
					// Выполняем логирование или отправляем сообщение о возникшей ошибке
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, что-то пошло не так. Пожалуйста, попробуйте снова."))
					log.Printf("Ошибка 400 при вызове CreateChatCompletion: %v\n", errorDetails)
					bot.Send(tgbotapi.NewMessage(435809098, errorDetails))
				} else {
					// Если ошибка не является ошибкой 400, обрабатываем ее соответствующим образом
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я устала запоминать, обнуляюсь"))
					log.Printf("Ошибка при вызове CreateChatCompletion: %v\n", err)
				}
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Message.Content)
			bot.Send(msg)
			req.Messages = append(req.Messages, resp.Choices[0].Message)
		}
		// Использует клиентскую функцию CreateImage для создания изображения на основе текстовой подсказки, предоставленной в аргументах команды.
		if update.Message.Command() == "imagine" {
			respUrl, err := client.CreateImage(
				context.Background(),
				openai.ImageRequest{
					Prompt:         update.Message.CommandArguments(),
					Size:           openai.CreateImageSize256x256,
					ResponseFormat: openai.CreateImageResponseFormatURL,
					N:              1,
				},
			)
			if err != nil {
				log.Printf("Image creation error: %v\n", err)
				continue
			}
			// Отправляет полученный URL изображения в чат с помощью Telegram API
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, respUrl.Data[0].URL)
			bot.Send(msg)
		}
	}
}
