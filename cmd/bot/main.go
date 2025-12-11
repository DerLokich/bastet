// cmd/bot/main.go
package main

import (
	"log"

	"BastetTetlegram/internal/bot"
	"BastetTetlegram/internal/config"
)

func main() {
	cfg := config.Load()

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	b.Run()
}
