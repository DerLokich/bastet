// internal/services/openai.go
package services

import (
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{client: client}
}

func (s *OpenAIService) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	// Добавьте таймаут, если нужно
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.CreateChatCompletion(ctx, req)
}

func (s *OpenAIService) CreateImage(ctx context.Context, req openai.ImageRequest) (openai.ImageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	return s.client.CreateImage(ctx, req)
}

// Добавьте другие методы, связанные с OpenAI
