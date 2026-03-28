package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	pool *pgxpool.Pool
}

func New(ctx context.Context, connectionString string) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 20
	poolCfg.MaxConnLifetime = 1 * time.Hour
	poolCfg.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), connectionString)
	return &DB{
		DB:   db,
		pool: pool,
	}, nil
}
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close(ctx context.Context) error {
	return db.DB.Close()
}
