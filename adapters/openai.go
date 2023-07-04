package adapters

import (
	"errors"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

var openAIAdapter *openai.Client

var (
	ErrOpenAIOrganizationIdEmpty  = errors.New("OpenAI organization id can't be empty")
	ErrOpenAIAPIKeyEmpty          = errors.New("OpenAI API key can't be empty")
	ErrOpenAIAdapterNotConfigured = errors.New("OpenAI adapter is not configured yet")
)

func ConfigureOpenAIAdapter(organizationId, apiKey string) {
	if organizationId == "" {
		log.Fatal().Err(ErrOpenAIOrganizationIdEmpty).Msg("Failed to configure OpenAI adapter")
	}
	if apiKey == "" {
		log.Fatal().Err(ErrOpenAIAPIKeyEmpty).Msg("Failed to configure OpenAI adapter")
	}

	config := openai.DefaultConfig(apiKey)
	config.OrgID = organizationId

	openAIAdapter = openai.NewClientWithConfig(config)
}

func OpenAIAdapter() *openai.Client {
	if openAIAdapter == nil {
		log.Fatal().Err(ErrOpenAIAdapterNotConfigured).Msg("Failed to get OpenAI adapter")
	}
	return openAIAdapter
}
