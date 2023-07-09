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

func generateSimplifiedText(ctx context.Context, originalText string) (simplifiedText string, err error) {
	systemPrompt := `Kamu bisa menjelaskan suatu topik yang kompleks dengan baik dan dapat membentuk penjelasan yang mudah dipahami orang. Tugasmu adalah untuk menyederhanakan teks yang akan diberikan menjadi bentuk yang lebih sederhana dan mudah dipahami. Kamu bebas mengurangi kata dan menggunakan bahasa yang lebih mudah jika perlu selama inti dari teksnya tetap tersampaikan.`
	prompt := fmt.Sprintf(`Saya kurang mengerti mengenai teks di bawah ini. Tolong disederhanakan agar saya bisa memahaminya dengan lebih mudah:

%s`, originalText)

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
			MaxTokens:   5000,
			Temperature: 0.8,
		},
	)

	if err != nil {
		switch e := err.(type) {
		case *openai.APIError:
			switch e.HTTPStatusCode {
			case http.StatusUnauthorized:
				return simplifiedText, ErrInvalidOpenAIAPIKey
			case http.StatusTooManyRequests:
				return simplifiedText, ErrOpenAIRateLimited
			case http.StatusInternalServerError, http.StatusServiceUnavailable:
				log.Err(ErrOpenAIServiceError).Msg("Failed to generate simplified text")
				return simplifiedText, ErrOpenAIServiceError
			default:
				log.Err(err).Msg("Failed to generate simplified text")
				return
			}
		default:
			log.Err(err).Msg("Failed to generate simplified text")
			return
		}
	}

	log.Info().Fields(map[string]any{
		"id":      res.ID,
		"model":   res.Model,
		"usage":   res.Usage,
		"choices": res.Choices,
	}).Msg("OpenAI - Generate Simplified Text Request")

	return res.Choices[0].Message.Content, nil
}

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
