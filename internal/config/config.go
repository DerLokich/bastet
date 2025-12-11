// internal/config/config.go
package config

import (
	"time"
)

type Config struct {
	Telegram struct {
		Token string `mapstructure:"token"`
	}
	OpenAI struct {
		Token string `mapstructure:"token"`
	}
	Files struct {
		PhrasesFile string `mapstructure:"phrases_file"`
		ToastsFile  string `mapstructure:"toasts_file"`
	}
	Storage struct {
		LastMentionFile string `mapstructure:"last_mention_file"`
	}
	Commands struct {
		Timeout time.Duration `mapstructure:"timeout"`
	}
	// Добавьте другие настройки
}
