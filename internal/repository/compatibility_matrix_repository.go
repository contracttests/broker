package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertCompatibilityMatrixQuery = `
		INSERT INTO compatibility_matrix
			(participant_id, version, counterpart_participant_id, counterpart_version, deployable)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	anyFailureSinceQuery = `
		SELECT EXISTS(
			SELECT 1 FROM compatibility_matrix
			WHERE participant_id = $1
			  AND version = $2
			  AND created_at >= $3
			  AND deployable = false
		)
	`
)

type CompatibilityMatrixRepository struct {
	pool *pgxpool.Pool
}

func NewCompatibilityMatrixRepository(pool *pgxpool.Pool) *CompatibilityMatrixRepository {
	return &CompatibilityMatrixRepository{pool: pool}
}

func (r *CompatibilityMatrixRepository) BeginTx(ctx context.Context) pgx.Tx {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(fmt.Errorf("error starting compatibility_matrix transaction: %w", err))
	}
	return tx
}

func (r *CompatibilityMatrixRepository) Insert(ctx context.Context, tx pgx.Tx, row *model.CompatibilityMatrixRow) {
	var counterpartID sql.NullInt64
	if row.CounterpartParticipantID != nil {
		counterpartID = sql.NullInt64{Int64: *row.CounterpartParticipantID, Valid: true}
	}

	var counterpartVersion sql.NullString
	if row.CounterpartVersion != nil {
		counterpartVersion = sql.NullString{String: *row.CounterpartVersion, Valid: true}
	}

	if err := tx.QueryRow(
		ctx,
		insertCompatibilityMatrixQuery,
		row.ParticipantID,
		row.Version,
		counterpartID,
		counterpartVersion,
		row.Deployable,
	).Scan(&row.ID, &row.CreatedAt); err != nil {
		panic(fmt.Errorf("error inserting compatibility matrix row: %w", err))
	}
}

func (r *CompatibilityMatrixRepository) AnyFailureSince(ctx context.Context, participantID int64, version string, since time.Time) bool {
	var exists bool
	if err := r.pool.QueryRow(ctx, anyFailureSinceQuery, participantID, version, since).Scan(&exists); err != nil {
		panic(fmt.Errorf("error checking compatibility matrix failures: %w", err))
	}
	return exists
}
