package article

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

func generateArticleText(ctx context.Context, originalDifficulty, targetDifficulty, text string) (generatedText string, err error) {
	systemPrompt := `Kamu bertugas untuk menyederhanakan bacaan sesuai dengan level pemahaman baca yang diinginkan. Ada tiga level pemahaman baca:

1. ADVANCED, ditujukan untuk teks yang butuh pemahaman baca tinggi. Seperti untuk orang-orang di dunia kerja dan mahasiswa.
2. INTERMEDIATE, ditujukan untuk teks yang butuh pemahaman baca menengah. Seperti siswa-siswa SMP kelas 7 di Indonesia sampai SMA kelas 12.
3. BEGINNER, ditujukan untuk teks yang butuh pemahaman baca pemula. Seperti siswa-siswa SD di Indonesia kelas 1 sampai 6.

  User akan memberi tahu kamu apa level pemahaman baca dari bacaan yang diberi serta level pemahaman baca yang user inginkan. Lalu di bawahnya, user akan memberikan bacaan yang akan kamu sederhanakan ke level pemahaman baca yang user inginkan
`
	prompt := fmt.Sprintf(`Teks di bawah ini dalam level pemahaman baca %s. Saya ingin kamu menyederhanakan teks berikut ke level pemahaman baca %s:

%s`, originalDifficulty, targetDifficulty, text)

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
			MaxTokens:   8000,
			Temperature: 0.8,
		},
	)

	if err != nil {
		switch e := err.(type) {
		case *openai.APIError:
			switch e.HTTPStatusCode {
			case http.StatusUnauthorized:
				return generatedText, ErrInvalidOpenAIAPIKey
			case http.StatusTooManyRequests:
				return generatedText, ErrOpenAIRateLimited
			case http.StatusInternalServerError, http.StatusServiceUnavailable:
				log.Err(ErrOpenAIServiceError).Msg("Failed to generate OpenAI article text")
				return generatedText, ErrOpenAIServiceError
			default:
				log.Err(err).Msg("Failed to generate OpenAI article text")
				return
			}
		default:
			log.Err(err).Msg("Failed to generate OpenAI article text")
			return
		}
	}

	log.Info().Fields(map[string]any{
		"id":      res.ID,
		"model":   res.Model,
		"usage":   res.Usage,
		"choices": res.Choices,
	}).Msg("OpenAI - Generate Article Text Request")

	return res.Choices[0].Message.Content, nil
}
