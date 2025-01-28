package main

import (
	"BastetTetlegram/config"
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const substr = "сосед"

const (
	cmdMe       = "me"
	cmdIDDQD    = "iddqd"
	cmdGPT      = "gpt"
	cmdImagine  = "imagine"
	cmdClaude   = "claude"
	cmdStart    = "start"
	cmdDeepSeek = "ds"
)

var titles = []string{"день", "дня", "дней"}

func main() {
	//Создается экземпляр бота, используя токен, полученный из config.Token
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

	ClaudeClient := anthropic.NewClient(option.WithAPIKey(config.ClaudeToken))

	LastMention := time.Now()

	DSClient := deepseek.NewClient(os.Getenv(config.DSToken))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue // Ignore any non-Message or non-command updates
		}
		messageText := update.Message.Text
		switch update.Message.Command() {
		// Данный фрагмент кода проверяет, является ли полученная команда от пользователя "me"
		case cmdMe:
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			bot.Request(kill)
		// Этот фрагмент кода позволяет боту устанавливать определенные права доступа для указанного пользователя в чате при получении команды "iddqd"
		case cmdStart:
			text := "normal text \n" +
				"*bold text* \n" +
				"_italic text_"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		case cmdIDDQD:
			// Создается переменная, которая используется для установки прав доступа для определенного пользователя в чате
			memberConfig := tgbotapi.PromoteChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: -1001165249098,
					UserID: 435809098,
				},
				// Устанавливается значение true для разрешения выполнения соответствующих действий
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
		case cmdGPT:
			ctx := context.Background()
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			})
			// Создаем канал для отмены контекста (если понадобится)
			ctx, cancel := context.WithCancel(ctx)
			resp, err := client.CreateChatCompletion(ctx, req)
			if err != nil {
				apiErr, ok := err.(*openai.APIError)
				// Если ошибка является ошибкой 400, выполняются следующие действия: отменяется контекст, обновляется
				if ok && apiErr.HTTPStatusCode == 400 {
					cancel()
					req = openai.ChatCompletionRequest{
						Model: openai.GPT4o,
						Messages: []openai.ChatCompletionMessage{
							{
								Role:    openai.ChatMessageRoleSystem,
								Content: "Clear message",
							},
						},
					}
					errorDetails := apiErr.Error()
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
		case cmdDeepSeek:
			ctx := context.Background()
			messages := []deepseek.ChatCompletionMessage{{
				Role:    constants.ChatMessageRoleUser,
				Content: update.Message.CommandArguments(),
			}}
			log.Printf("Вызываем ключ %v", config.DSToken)
			response1, err := DSClient.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
				Model:    deepseek.DeepSeekChat,
				Messages: messages,
			})
			if err != nil {
				log.Fatalf("Round 1 failed: %v", err)
			}

			messages = append(messages, deepseek.ChatCompletionMessage{
				Role:    response1.Choices[0].Message.Role,
				Content: response1.Choices[0].Message.Content,
			})
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response1.Choices[0].Message.Content)
			bot.Send(msg)

			log.Printf("Messages after Round 1: %+v\n", messages)
			//
			//messages = append(messages, deepseek.ChatCompletionMessage{
			//	Role:    constants.ChatMessageRoleUser,
			//	Content: "What is the second?",
			//})
			//
			//response2, err := DSCLient.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
			//	Model:    deepseek.DeepSeekChat,
			//	Messages: messages,
			//})
			//if err != nil {
			//	log.Fatalf("Round 2 failed: %v", err)
			//}
			//
			//fmt.Printf("Final messages: %+v\n", append(messages, deepseek.ChatCompletionMessage{
			//	Role:    response2.Choices[0].Message.Role,
			//	Content: response2.Choices[0].Message.Content,
			//}))
		case cmdClaude:
			response, err := ClaudeClient.Messages.New(context.TODO(), anthropic.MessageNewParams{
				Model:     anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.MessageParam{
					anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
				}),
			})
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.Content[0].Text)
			bot.Send(msg)

		// Использует клиентскую функцию CreateImage для создания изображения на основе текстовой подсказки, предоставленной в аргументах команды
		case cmdImagine:
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
		default:
			// Ignore any unrecognized commands
		}

		// Проверяет, содержит ли текст сообщения подстроку

		if strings.Contains(strings.ToLower(messageText), substr) {
			// Вычисляет разницу времени с момента последнего упоминания в днях
			TimeDifference := time.Since(LastMention).Hours() / 24
			// Создает сообщение с текстом, содержащим полученную разницу времени и отправляет его в чат
			Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Neib))
			log.Println(TimeDifference)
			log.Printf(LastMention.String())
			LastMention = time.Now()
			log.Printf(LastMention.String())
		}
	}

}

// declOfNum returns the proper form of a noun based on the given number.
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
