// internal/config/load.go
package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	// Также можно загружать из переменных окружения
	viper.AutomaticEnv()

	// Установим значения по умолчанию
	viper.SetDefault("files.phrases_file", "config/phrases.txt")
	viper.SetDefault("files.toasts_file", "config/toasts.txt")
	viper.SetDefault("storage.last_mention_file", "last_mention.json")
	viper.SetDefault("commands.timeout", 30*time.Second)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	// Загрузка токенов из переменных окружения
	cfg.Telegram.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	if cfg.Telegram.Token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	cfg.OpenAI.Token = os.Getenv("OPENAI_API_KEY")
	if cfg.OpenAI.Token == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	return &cfg
}
