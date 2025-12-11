// internal/models/models.go
package models

import "time"

type LastMentionData struct {
	LastMention time.Time `json:"last_mention"`
}
