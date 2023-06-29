package article

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	mediastackUrl    string
	mediastackApiKey string

	ErrEmptyMediastackUrl    = errors.New("Empty mediastack url")
	ErrEmptyMediastackApiKey = errors.New("Empty mediastack api key")
)

func ConfigureMediastackAdapter(url, apiKey string) {
	if strings.TrimSpace(url) == "" {
		log.Fatal().Err(ErrEmptyMediastackUrl).Msg("Failed to configure mediastack adapter")
	}
	if strings.TrimSpace(apiKey) == "" {
		log.Fatal().Err(ErrEmptyMediastackApiKey).Msg("Failed to configure mediastack adapter")
	}

	mediastackUrl = url
	mediastackApiKey = apiKey
}
