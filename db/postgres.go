package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func CreateConnPool(dsn string) (pool *pgxpool.Pool) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create connection pool")
	}

	_, err = pool.Acquire(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Info().Msg("Established connection to database!")
	return pool
}
