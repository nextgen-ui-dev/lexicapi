package adapters

import (
	"errors"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

var (
	ErrOpenAIOrganizationIdEmpty = errors.New("OpenAI organization id can't be empty")
	ErrOpenAIAPIKeyEmpty         = errors.New("OpenAI API key can't be empty")
)

func ConfigureOpenAIAdapter(organizationId, apiKey string) *openai.Client {
	if organizationId == "" {
		log.Fatal().Err(ErrOpenAIOrganizationIdEmpty).Msg("Failed to configure OpenAI adapter")
	}
	if apiKey == "" {
		log.Fatal().Err(ErrOpenAIAPIKeyEmpty).Msg("Failed to configure OpenAI adapter")
	}

	config := openai.DefaultConfig(apiKey)
	config.OrgID = organizationId

	return openai.NewClientWithConfig(config)
}
