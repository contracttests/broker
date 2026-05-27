package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClockRepository struct {
	pool *pgxpool.Pool
}

func NewClockRepository(pool *pgxpool.Pool) *ClockRepository {
	return &ClockRepository{pool: pool}
}

func (r *ClockRepository) DBNow(ctx context.Context) time.Time {
	var now time.Time
	if err := r.pool.QueryRow(ctx, `SELECT now()`).Scan(&now); err != nil {
		panic(fmt.Errorf("error reading db clock: %w", err))
	}
	return now
}
