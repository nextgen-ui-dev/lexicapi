package article

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

var (
	openAIAdapter *openai.Client

	ErrNilOpenAIAdapter = errors.New("OpenAI adapter can't be nil")
)

func SetOpenAIAdapter(adapter *openai.Client) {
	if adapter == nil {
		log.Fatal().Err(ErrNilOpenAIAdapter).Msg("Failed to set OpenAI adapter for article module")
	}

	openAIAdapter = adapter
}
