package friend

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	pool *pgxpool.Pool

	ErrNilPool = errors.New("Connection pool can't be nil")
)

func SetPool(newPool *pgxpool.Pool) {
	if newPool == nil {
		log.Fatal().Err(ErrNilPool).Msg("Failed to set connection pool for friend module")
	}

	pool = newPool
}
