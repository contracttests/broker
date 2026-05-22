package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/wiredb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProviderResourceNotFound = errors.New("provider resource not found")
)

const (
	insertContractQuery = `
		INSERT INTO contracts
			(uuid, name, owner)
		VALUES
			($1, $2, $3)
		RETURNING id
	`

	insertContractVersionQuery = `
		INSERT INTO contract_versions
			(uuid, contract_id, version, checksum, raw_payload)
		SELECT
			$1, $2, COALESCE(MAX(version), 0) + 1, $3, $4
		FROM
			contract_versions
		WHERE
			contract_id = $2
		RETURNING id, version
	`

	insertResourceQuery = `
		INSERT INTO resources
			(
				uuid,
				contract_id,
				direction,
				kind,
				provider,
				endpoint,
				method,
				status_code,
				provider_hash,
				consumer_hash
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	insertPropertyQuery = `
		INSERT INTO properties
			(uuid, resource_id, path)
		VALUES
			($1, $2, $3)
		RETURNING id
	`

	insertPropertyVersionQuery = `
		INSERT INTO property_versions
			(uuid, property_id, contract_version_id, type, optional, change)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	findContractTreeQuery = `
		SELECT
			c.id,
			c.uuid,
			c.name,
			c.owner,
			c.created_at,
			r.id,
			r.direction,
			r.kind,
			r.provider,
			r.endpoint,
			r.method,
			r.status_code,
			r.provider_hash,
			r.consumer_hash,
			r.created_at,
			p.id,
			p.path,
			pv.type,
			pv.optional,
			pv.change
		FROM contracts c
		JOIN resources r ON r.contract_id = c.id
		JOIN properties p ON p.resource_id = r.id
		JOIN LATERAL (
			SELECT type, optional, change
			FROM property_versions
			WHERE property_id = p.id
			ORDER BY contract_version_id DESC
			LIMIT 1
		) pv ON true
		WHERE c.name = $1
		ORDER BY r.id`

	findResourcesByDirectionAndProviderHashQuery = `
		SELECT
			r.id,
			r.uuid,
			r.contract_id,
			c.name,
			c.owner,
			r.direction,
			r.kind,
			r.provider,
			r.endpoint,
			r.method,
			r.status_code,
			r.provider_hash,
			r.consumer_hash,
			r.created_at,
			p.path,
			pv.type,
			pv.optional,
			pv.change
		FROM
			contracts c
		JOIN
			resources r ON r.contract_id = c.id
		LEFT JOIN
			properties p ON p.resource_id = r.id
		LEFT JOIN LATERAL (
			SELECT
				type,
				optional,
				change
			FROM
				property_versions
			WHERE
				property_id = p.id
			ORDER BY contract_version_id DESC
			LIMIT 1
		) pv ON true
		WHERE
			r.direction = $1
		AND
			r.provider_hash = $2
		AND
			(pv.change IS NULL OR pv.change != 'removed')`
)

type ContractRepository struct {
	pool *pgxpool.Pool
}

func NewContractRepository(pool *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{pool: pool}
}

func (r *ContractRepository) ExistsByName(ctx context.Context, name string) bool {
	const query = `SELECT EXISTS(SELECT 1 FROM contracts WHERE name = $1)`

	var exists bool

	err := r.pool.QueryRow(
		ctx,
		query,
		name,
	).Scan(
		&exists,
	)
	if err != nil {
		panic(fmt.Errorf("error finding contract by name: %w", err))
	}

	return exists
}

func (r *ContractRepository) Save(
	ctx context.Context,
	contract *model.Contract,
) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(fmt.Errorf("error starting transaction: %w", err))
	}

	defer tx.Rollback(ctx)

	contractRow := wiredb.NewInsertContractRow(contract)
	r.insertContract(ctx, tx, contractRow)
	contract.ID = contractRow.ID

	contractVersion := model.NewContractVersion(contract)
	contractVersionRow := wiredb.NewInsertContractVersionRow(contractVersion)
	r.insertContractVersion(ctx, tx, contractVersionRow)
	contractVersion.ID = contractVersionRow.ID

	for _, resource := range contract.Resources {
		resourceRow := wiredb.NewInsertResourceRow(contract, resource)
		r.insertResource(ctx, tx, resourceRow)
		resource.ID = resourceRow.ID

		for _, property := range resource.Properties {
			propertyRow := wiredb.NewInsertPropertyRow(resource, property)
			r.insertNewProperty(ctx, tx, propertyRow)
			property.ID = propertyRow.ID

			propertyVersionRow := wiredb.NewInsertPropertyVersionRowAdded(contractVersion, property)
			r.insertPropertyVersion(ctx, tx, propertyVersionRow)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(fmt.Errorf("error committing transaction: %w", err))
	}
}

func (r *ContractRepository) Update(
	ctx context.Context,
	next *model.Contract,
) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(fmt.Errorf("error starting transaction: %w", err))
	}

	defer tx.Rollback(ctx)

	contract := r.LoadLatestContractByName(ctx, next.Name)

	diff := contract.Diff(next)

	if len(diff.Resources) == 0 {
		if err := tx.Commit(ctx); err != nil {
			panic(fmt.Errorf("error committing transaction: %w", err))
		}

		return
	}

	next.ID = contract.ID

	contractVersion := model.NewContractVersion(next)
	contractVersionRow := wiredb.NewInsertContractVersionRow(contractVersion)
	r.insertContractVersion(ctx, tx, contractVersionRow)
	contractVersion.ID = contractVersionRow.ID

	for _, resourceChange := range diff.Resources {
		switch resourceChange.Kind {
		case model.ChangeAdded:
			resource := resourceChange.Resource
			resourceRow := wiredb.NewInsertResourceRow(next, resource)
			r.insertResource(ctx, tx, resourceRow)
			resource.ID = resourceRow.ID

			for _, property := range resource.Properties {
				propertyRow := wiredb.NewInsertPropertyRow(resource, property)
				r.insertNewProperty(ctx, tx, propertyRow)
				property.ID = propertyRow.ID

				propertyVersionRow := wiredb.NewInsertPropertyVersionRowAdded(contractVersion, property)
				r.insertPropertyVersion(ctx, tx, propertyVersionRow)
			}

		case model.ChangeModified:
			resource := contract.Resources[resourceChange.Resource.PrimaryHash()]

			for _, propertyChange := range resourceChange.Properties {
				switch propertyChange.Kind {
				case model.ChangeAdded:
					property := propertyChange.After
					propertyRow := wiredb.NewInsertPropertyRow(resource, property)
					r.insertNewProperty(ctx, tx, propertyRow)
					property.ID = propertyRow.ID

					propertyVersionRow := wiredb.NewInsertPropertyVersionRowAdded(contractVersion, property)
					r.insertPropertyVersion(ctx, tx, propertyVersionRow)

				case model.ChangeModified:
					property := propertyChange.After
					property.ID = resource.Properties[propertyChange.After.Path].ID

					propertyVersionRow := wiredb.NewInsertPropertyVersionRowModified(contractVersion, property)
					r.insertPropertyVersion(ctx, tx, propertyVersionRow)

				case model.ChangeRemoved:
					property := resource.Properties[propertyChange.Before.Path]

					propertyVersionRow := wiredb.NewInsertPropertyVersionRowRemoved(contractVersion, property)
					r.insertPropertyVersion(ctx, tx, propertyVersionRow)
				}
			}

		case model.ChangeRemoved:
			resource := contract.Resources[resourceChange.Resource.PrimaryHash()]

			for _, propertyChange := range resourceChange.Properties {
				property := resource.Properties[propertyChange.Before.Path]

				propertyVersionRow := wiredb.NewInsertPropertyVersionRowRemoved(contractVersion, property)
				r.insertPropertyVersion(ctx, tx, propertyVersionRow)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(fmt.Errorf("error committing transaction: %w", err))
	}
}

func (r *ContractRepository) insertContract(
	ctx context.Context,
	tx pgx.Tx,
	row *wiredb.InsertContractRow,
) {
	if err := tx.QueryRow(
		ctx,
		insertContractQuery,
		row.UUID,
		row.Name,
		row.Owner,
	).Scan(&row.ID); err != nil {
		panic(fmt.Errorf("error inserting contract: %w", err))
	}
}

func (r *ContractRepository) insertContractVersion(
	ctx context.Context,
	tx pgx.Tx,
	row *wiredb.InsertContractVersionRow,
) {
	if err := tx.QueryRow(
		ctx,
		insertContractVersionQuery,
		row.UUID,
		row.ContractID,
		row.Checksum,
		row.RawPayload,
	).Scan(&row.ID, &row.Version); err != nil {
		panic(fmt.Errorf("error inserting contract version: %w", err))
	}
}

func (r *ContractRepository) insertResource(
	ctx context.Context,
	tx pgx.Tx,
	row *wiredb.InsertResourceRow,
) {
	if err := tx.QueryRow(
		ctx,
		insertResourceQuery,
		row.UUID,
		row.ContractID,
		row.Direction,
		row.Kind,
		row.Provider,
		row.Endpoint,
		row.Method,
		row.StatusCode,
		row.ProviderHash,
		row.ConsumerHash,
	).Scan(&row.ID); err != nil {
		panic(fmt.Errorf("error inserting resource: %w", err))
	}
}

func (r *ContractRepository) insertNewProperty(
	ctx context.Context,
	tx pgx.Tx,
	row *wiredb.InsertPropertyRow,
) {
	if err := tx.QueryRow(
		ctx,
		insertPropertyQuery,
		row.UUID,
		row.ResourceID,
		row.Path,
	).Scan(&row.ID); err != nil {
		panic(fmt.Errorf("error inserting property: %w", err))
	}
}

func (r *ContractRepository) insertPropertyVersion(
	ctx context.Context,
	tx pgx.Tx,
	row *wiredb.InsertPropertyVersionRow,
) {
	if _, err := tx.Exec(
		ctx,
		insertPropertyVersionQuery,
		row.UUID,
		row.PropertyID,
		row.ContractVersionID,
		row.Type,
		row.Optional,
		row.Change,
	); err != nil {
		panic(fmt.Errorf("error inserting property version: %w", err))
	}
}

func (r *ContractRepository) LoadLatestContractByName(
	ctx context.Context,
	contractName string,
) *model.Contract {
	rows, err := r.pool.Query(ctx, findContractTreeQuery, contractName)

	if err != nil {
		panic(fmt.Errorf("error loading contract tree: %w", err))
	}

	defer rows.Close()

	var found bool
	var contract *model.Contract

	for rows.Next() {
		var row wiredb.TableRow

		if err := rows.Scan(
			&row.ContractID,
			&row.ContractUUID,
			&row.ContractName,
			&row.ContractOwner,
			&row.ContractCreatedAt,
			&row.ResourceID,
			&row.ResourceDirection,
			&row.ResourceKind,
			&row.ResourceProvider,
			&row.ResourceEndpoint,
			&row.ResourceMethod,
			&row.ResourceStatusCode,
			&row.ResourceProviderHash,
			&row.ResourceConsumerHash,
			&row.ResourceCreatedAt,
			&row.PropertyID,
			&row.PropertyPath,
			&row.PropertyVersionType,
			&row.PropertyVersionOptional,
			&row.PropertyVersionChange,
		); err != nil {
			panic(fmt.Errorf("error scanning contract tree row: %w", err))
		}

		if row.PropertyVersionChange == string(model.ChangeRemoved) {
			continue
		}

		if !found {
			contract = row.ToContractModel()
			found = true
		}

		resource := row.ToResourceModel()
		key := resource.PrimaryHash()
		if _, seen := contract.Resources[key]; !seen {
			contract.AddResource(resource)
		}

		property := row.ToPropertyModel()
		if _, seen := contract.Resources[key].Properties[property.Path]; !seen {
			contract.Resources[key].Properties[property.Path] = property
		}
	}

	return contract
}

func (r *ContractRepository) LoadProviderResource(ctx context.Context, consumer model.Resource) (model.Resource, error) {
	rows, err := r.pool.Query(
		ctx,
		findResourcesByDirectionAndProviderHashQuery,
		string(model.Provides),
		consumer.ProviderHash(),
	)

	if err != nil {
		return model.Resource{}, ErrProviderResourceNotFound
	}

	defer rows.Close()

	var found bool
	var provider model.Resource

	for rows.Next() {
		var row wiredb.TableRow

		if err := rows.Scan(
			&row.ResourceID,
			&row.ResourceUUID,
			&row.ContractID,
			&row.ContractName,
			&row.ContractOwner,
			&row.ResourceDirection,
			&row.ResourceKind,
			&row.ResourceProvider,
			&row.ResourceEndpoint,
			&row.ResourceMethod,
			&row.ResourceStatusCode,
			&row.ResourceProviderHash,
			&row.ResourceConsumerHash,
			&row.ResourceCreatedAt,
			&row.PropertyPath,
			&row.PropertyVersionType,
			&row.PropertyVersionOptional,
			&row.PropertyVersionChange,
		); err != nil {
			panic(fmt.Errorf("error scanning provider resource: %w", err))
		}

		if !found {
			provider = row.ToResourceModel()
			found = true
		}

		provider.Properties[row.PropertyPath] = row.ToPropertyModel()
	}

	if !found {
		return model.Resource{}, ErrProviderResourceNotFound
	}

	return provider, nil
}

func (r *ContractRepository) FindConsumersOfProvider(ctx context.Context, provider model.Resource) []model.Resource {
	rows, err := r.pool.Query(
		ctx,
		findResourcesByDirectionAndProviderHashQuery,
		string(model.Consumes),
		provider.ProviderHash(),
	)

	if err != nil {
		panic(fmt.Errorf("error finding consumers of provider: %w", err))
	}

	defer rows.Close()

	consumersMap := make(map[int64]model.Resource)

	for rows.Next() {
		var row wiredb.TableRow

		if err := rows.Scan(
			&row.ResourceID,
			&row.ResourceUUID,
			&row.ContractID,
			&row.ContractName,
			&row.ContractOwner,
			&row.ResourceDirection,
			&row.ResourceKind,
			&row.ResourceProvider,
			&row.ResourceEndpoint,
			&row.ResourceMethod,
			&row.ResourceStatusCode,
			&row.ResourceProviderHash,
			&row.ResourceConsumerHash,
			&row.ResourceCreatedAt,
			&row.PropertyPath,
			&row.PropertyVersionType,
			&row.PropertyVersionOptional,
			&row.PropertyVersionChange,
		); err != nil {
			panic(fmt.Errorf("error scanning consumer: %w", err))
		}

		if _, seen := consumersMap[row.ResourceID]; !seen {
			consumersMap[row.ResourceID] = row.ToResourceModel()
		}

		consumersMap[row.ResourceID].Properties[row.PropertyPath] = row.ToPropertyModel()
	}

	consumers := make([]model.Resource, 0, len(consumersMap))

	for _, consumer := range consumersMap {
		consumers = append(consumers, consumer)
	}

	return consumers
}
