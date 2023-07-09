package assistant

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

var (
	ErrInvalidOpenAIAPIKey = errors.New("Invalid OpenAI API key")
	ErrOpenAIRateLimited   = errors.New("OpenAI has rate limited us due to too many requests. Please try again later.")
	ErrOpenAIServiceError  = errors.New("OpenAI service is currently unavailable. Please try again later")
)

func generateTextExplanation(ctx context.Context, text string) (explanation string, err error) {
	systemPrompt := `Kamu adalah seorang pakar yang ahli dalam berbagai macam bidang dan pengetahuan yang kamu miliki luas. Tugas kamu adalah menjelaskan kalimat, paragraf, atau teks yang akan diberikan`
	prompt := fmt.Sprintf(`Tolong berikan penjelasan yang mudah dipahami mengenai teks berikut:

"%s"`, text)

	res, err := openAIAdapter.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   2000,
			Temperature: 0.8,
		},
	)

	if err != nil {
		switch e := err.(type) {
		case *openai.APIError:
			switch e.HTTPStatusCode {
			case http.StatusUnauthorized:
				return explanation, ErrInvalidOpenAIAPIKey
			case http.StatusTooManyRequests:
				return explanation, ErrOpenAIRateLimited
			case http.StatusInternalServerError, http.StatusServiceUnavailable:
				log.Err(ErrOpenAIServiceError).Msg("Failed to generate text explanation")
				return explanation, ErrOpenAIServiceError
			default:
				log.Err(err).Msg("Failed to generate text explanation")
				return
			}
		default:
			log.Err(err).Msg("Failed to generate text explanation")
			return
		}
	}

	log.Info().Fields(map[string]any{
		"id":      res.ID,
		"model":   res.Model,
		"usage":   res.Usage,
		"choices": res.Choices,
	}).Msg("OpenAI - Generate Explanation Text Request")

	return res.Choices[0].Message.Content, nil
}
