package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/contracttests/broker/server/pkg/dbhelper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FoundContract struct {
	ID    int64
	UUID  uuid.UUID
	Name  string
	Owner string
}

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}

func (r *ContractRepository) FindByName(ctx context.Context, name string) (*FoundContract, error) {
	const query = `SELECT id, uuid, name, owner FROM contracts WHERE name = $1`

	var found FoundContract
	err := r.pool.QueryRow(ctx, query, name).Scan(&found.ID, &found.UUID, &found.Name, &found.Owner)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error finding contract by name: %w", err)
	}

	return &found, nil
}

func (r *ContractRepository) Save(ctx context.Context, c *model.Contract, rawPayload []byte) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	insertContractQuery := `
		INSERT INTO contracts 
			(uuid, name, owner)
		VALUES 
			($1, $2, $3)
		RETURNING id
	`

	var contractID int64
	if err := tx.QueryRow(
		ctx,
		insertContractQuery,
		uuid.New(),
		c.Name,
		c.Owner,
	).Scan(&contractID); err != nil {
		return fmt.Errorf("error inserting contract: %w", err)
	}

	insertContractVersionQuery := `
		INSERT INTO contract_versions 
			(uuid, contract_id, version, checksum, raw_payload)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING id
	`

	var contractVersionID int64
	if err := tx.QueryRow(
		ctx,
		insertContractVersionQuery,
		uuid.New(),
		contractID,
		1,
		c.Checksum(),
		rawPayload,
	).Scan(&contractVersionID); err != nil {
		return fmt.Errorf("error inserting contract version: %w", err)
	}

	insertResourceQuery := `
		INSERT INTO resources
			(uuid, contract_id, direction, kind, provider, endpoint, method, status_code)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	insertPropertyQuery := `
		INSERT INTO properties 
			(uuid, resource_id, path) 
		VALUES 
			($1, $2, $3) 
		RETURNING id
	`

	insertPropertyVersionQuery := `
		INSERT INTO property_versions 
			(uuid, property_id, contract_version_id, type, optional, change) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, resource := range c.Resources {
		var resourceID int64
		if err := tx.QueryRow(
			ctx,
			insertResourceQuery,
			uuid.New(),
			contractID,
			string(resource.Direction),
			string(resource.Kind),
			dbhelper.NullableString(resource.Provider),
			resource.Endpoint,
			resource.Method,
			dbhelper.NullableString(resource.StatusCode),
		).Scan(&resourceID); err != nil {
			return fmt.Errorf("error inserting resource: %w", err)
		}

		for _, property := range resource.Properties {
			var propertyID int64
			if err := tx.QueryRow(
				ctx,
				insertPropertyQuery,
				uuid.New(),
				resourceID,
				property.Path,
			).Scan(&propertyID); err != nil {
				return fmt.Errorf("error inserting property: %w", err)
			}

			if _, err := tx.Exec(
				ctx,
				insertPropertyVersionQuery,
				uuid.New(),
				propertyID,
				contractVersionID,
				property.Type,
				property.Optional,
				"added",
			); err != nil {
				return fmt.Errorf("error inserting property version: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
