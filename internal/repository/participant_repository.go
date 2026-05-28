package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/contracttesting/broker/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	existsParticipantByNameQuery = `
		SELECT EXISTS(
			SELECT
				1
			FROM 
				participants
			WHERE 
				name = $1
		)
	`

	findParticipantByNameQuery = `
		SELECT
			id, name
		FROM
			participants
		WHERE
			name = $1
	`

	insertParticipantQuery = `
		INSERT INTO participants
			(name)
		VALUES
			($1)
		RETURNING id
	`
)

type ParticipantRepository struct {
	pool *pgxpool.Pool
}

func NewParticipantRepository(pool *pgxpool.Pool) *ParticipantRepository {
	return &ParticipantRepository{pool: pool}
}

func (r *ParticipantRepository) ExistsByName(ctx context.Context, name string) bool {
	var exists bool
	if err := r.pool.QueryRow(ctx, existsParticipantByNameQuery, name).Scan(&exists); err != nil {
		panic(fmt.Errorf("error finding participant by name: %w", err))
	}
	return exists
}

func (r *ParticipantRepository) FindByName(ctx context.Context, name string) (*model.Participant, bool) {
	p := &model.Participant{}
	err := r.pool.QueryRow(ctx, findParticipantByNameQuery, name).Scan(&p.ID, &p.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false
	}
	if err != nil {
		panic(fmt.Errorf("error finding participant by name: %w", err))
	}
	return p, true
}

func (r *ParticipantRepository) Create(ctx context.Context, p *model.Participant) {
	if err := r.pool.QueryRow(ctx, insertParticipantQuery, p.Name).Scan(&p.ID); err != nil {
		panic(fmt.Errorf("error inserting participant: %w", err))
	}
}
