package infrastructures

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgx() (libpgx.RDBMS, libpgx.Tx, squirrel.StatementBuilderType, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := os.Getenv("DATABASE_URI")
	if connString == "" {
		return nil, nil, squirrel.StatementBuilder, nil, fmt.Errorf("DATABASE_URI is not set")
	}

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, nil, squirrel.StatementBuilder, nil, fmt.Errorf("parse connection config: %w", err)
	}

	otelTracer := otelpgx.NewTracer()
	cfg.ConnConfig.Tracer = otelTracer

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, nil, squirrel.StatementBuilder, nil, fmt.Errorf("connect to database: %w", err)
	}

	rdbms := libpgx.NewRDBMS(pool)

	log.Println("inisiate pgx pool successfully")
	return rdbms, rdbms, squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar), func() {
		pool.Close()
	}, nil
}
