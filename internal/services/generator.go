// internal/services/generator.go
package services

import (
	"math/rand"
	"strings"
	"time"
)

// globalRand объявлен как переменная пакета
var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var festiveEmojis = []string{
	"🥂", "🍾", "🍷", "🍸", "🍺", "🍻", "🥂", "🎉", "🎊", "🥳", "✨", "🌟", "💫", "🔥", "❤️", "💖", "💕", "🌹", "💐", "🎁", "🎀", "🎊", "🎉", "🥂", "-toast-emoji-",
}

func GetRandomPhrase(phrases []string) string {
	if len(phrases) == 0 {
		return "Фразы закончились :("
	}
	return phrases[globalRand.Intn(len(phrases))]
}

func GetRandomToast(toasts []string) string {
	if len(toasts) == 0 {
		return "Тосты закончились :("
	}
	return toasts[globalRand.Intn(len(toasts))]
}

func GetRandomEmoji() string {
	if len(festiveEmojis) == 0 {
		return ""
	}
	return festiveEmojis[globalRand.Intn(len(festiveEmojis))]
}

func EscapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!", "[", "]", "(", ")"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
