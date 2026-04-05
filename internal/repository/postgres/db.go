package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepo envuelve el pool de conexiones
type PostgresRepo struct {
	Pool *pgxpool.Pool
}

// NewPostgresRepo configura y abre la conexión
func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Configuración Senior: No satures la DB, mantén un pool sano
	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error conectando a la DB: %w", err)
	}

	return &PostgresRepo{Pool: pool}, nil
}
