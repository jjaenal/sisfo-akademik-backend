package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DefaultConnTimeout = 10 * time.Second

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 10
	}
	cctx, cancel := context.WithTimeout(ctx, DefaultConnTimeout)
	defer cancel()
	p, err := pgxpool.NewWithConfig(cctx, cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

type txStarter interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error)
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func WithTx(ctx context.Context, starter txStarter, fn func(tx Tx) error) error {
	tx, err := starter.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
