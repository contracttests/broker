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

const (
	insertContractVersionQuery = `
		INSERT INTO contract_versions
			(uuid, contract_id, version, checksum, raw_payload)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id
	`

	insertResourceQuery = `
		INSERT INTO resources
			(uuid, contract_id, direction, kind, provider, endpoint, method, status_code)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
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

	findResourceQuery = `
		SELECT id FROM resources
		WHERE contract_id = $1
		  AND direction = $2
		  AND kind = $3
		  AND provider IS NOT DISTINCT FROM $4
		  AND endpoint = $5
		  AND method = $6
		  AND status_code IS NOT DISTINCT FROM $7
	`

	findPropertyQuery = `SELECT id FROM properties WHERE resource_id = $1 AND path = $2`
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

func (r *ContractRepository) FindByName(ctx context.Context, name string) *FoundContract {
	const query = `SELECT id, uuid, name, owner FROM contracts WHERE name = $1`

	var found FoundContract
	err := r.pool.QueryRow(ctx, query, name).Scan(&found.ID, &found.UUID, &found.Name, &found.Owner)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		panic(fmt.Errorf("error finding contract by name: %w", err))
	}

	return &found
}

func (r *ContractRepository) Save(ctx context.Context, c *model.Contract, rawPayload []byte) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(fmt.Errorf("error starting transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	const insertContractQuery = `
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
		panic(fmt.Errorf("error inserting contract: %w", err))
	}

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
		panic(fmt.Errorf("error inserting contract version: %w", err))
	}

	for _, resource := range c.Resources {
		resourceID := insertResource(ctx, tx, contractID, resource)
		for _, property := range resource.Properties {
			insertNewProperty(ctx, tx, resourceID, contractVersionID, property, model.ChangeAdded)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(fmt.Errorf("error committing transaction: %w", err))
	}
}

func (r *ContractRepository) Update(ctx context.Context, next *model.Contract, rawPayload []byte) {
	existing := r.FindByName(ctx, next.Name)
	if existing == nil {
		panic(fmt.Errorf("contract not found: %s", next.Name))
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		panic(fmt.Errorf("error starting transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	prev := loadLatestContract(ctx, tx, existing.ID, existing.Name, existing.Owner)

	diff := prev.Diff(next)
	if len(diff.Resources) == 0 {
		if err := tx.Commit(ctx); err != nil {
			panic(fmt.Errorf("error committing transaction: %w", err))
		}
		return
	}

	var nextVersion int
	if err := tx.QueryRow(
		ctx,
		`SELECT COALESCE(MAX(version), 0) + 1 FROM contract_versions WHERE contract_id = $1`,
		existing.ID,
	).Scan(&nextVersion); err != nil {
		panic(fmt.Errorf("error computing next version: %w", err))
	}

	var contractVersionID int64
	if err := tx.QueryRow(
		ctx,
		insertContractVersionQuery,
		uuid.New(),
		existing.ID,
		nextVersion,
		next.Checksum(),
		rawPayload,
	).Scan(&contractVersionID); err != nil {
		panic(fmt.Errorf("error inserting contract version: %w", err))
	}

	for _, rc := range diff.Resources {
		persistResourceChange(ctx, tx, existing.ID, contractVersionID, rc)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(fmt.Errorf("error committing transaction: %w", err))
	}
}

func persistResourceChange(ctx context.Context, tx pgx.Tx, contractID, contractVersionID int64, rc model.ResourceChange) {
	switch rc.Kind {
	case model.ChangeAdded:
		resourceID := insertResource(ctx, tx, contractID, rc.Resource)
		for _, property := range rc.Resource.Properties {
			insertNewProperty(ctx, tx, resourceID, contractVersionID, property, model.ChangeAdded)
		}

	case model.ChangeModified:
		resourceID := findResourceID(ctx, tx, contractID, rc.Resource)
		for _, pc := range rc.Properties {
			persistPropertyChange(ctx, tx, resourceID, contractVersionID, pc)
		}

	case model.ChangeRemoved:
		resourceID := findResourceID(ctx, tx, contractID, rc.Resource)
		for _, pc := range rc.Properties {
			propertyID := findPropertyID(ctx, tx, resourceID, pc.Before.Path)
			insertPropertyVersion(ctx, tx, propertyID, contractVersionID, pc.Before, model.ChangeRemoved)
		}
	}
}

func persistPropertyChange(ctx context.Context, tx pgx.Tx, resourceID, contractVersionID int64, pc model.PropertyChange) {
	switch pc.Kind {
	case model.ChangeAdded:
		insertNewProperty(ctx, tx, resourceID, contractVersionID, pc.After, model.ChangeAdded)
	case model.ChangeModified:
		propertyID := findPropertyID(ctx, tx, resourceID, pc.After.Path)
		insertPropertyVersion(ctx, tx, propertyID, contractVersionID, pc.After, model.ChangeModified)
	case model.ChangeRemoved:
		propertyID := findPropertyID(ctx, tx, resourceID, pc.Before.Path)
		insertPropertyVersion(ctx, tx, propertyID, contractVersionID, pc.Before, model.ChangeRemoved)
	}
}

func insertResource(ctx context.Context, tx pgx.Tx, contractID int64, resource model.Resource) int64 {
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
		panic(fmt.Errorf("error inserting resource: %w", err))
	}
	return resourceID
}

func insertNewProperty(ctx context.Context, tx pgx.Tx, resourceID, contractVersionID int64, property model.Property, change model.ChangeKind) {
	var propertyID int64
	if err := tx.QueryRow(
		ctx,
		insertPropertyQuery,
		uuid.New(),
		resourceID,
		property.Path,
	).Scan(&propertyID); err != nil {
		panic(fmt.Errorf("error inserting property: %w", err))
	}
	insertPropertyVersion(ctx, tx, propertyID, contractVersionID, property, change)
}

func insertPropertyVersion(ctx context.Context, tx pgx.Tx, propertyID, contractVersionID int64, property model.Property, change model.ChangeKind) {
	if _, err := tx.Exec(
		ctx,
		insertPropertyVersionQuery,
		uuid.New(),
		propertyID,
		contractVersionID,
		property.Type,
		property.Optional,
		string(change),
	); err != nil {
		panic(fmt.Errorf("error inserting property version: %w", err))
	}
}

func findResourceID(ctx context.Context, tx pgx.Tx, contractID int64, resource model.Resource) int64 {
	var resourceID int64
	if err := tx.QueryRow(
		ctx,
		findResourceQuery,
		contractID,
		string(resource.Direction),
		string(resource.Kind),
		dbhelper.NullableString(resource.Provider),
		resource.Endpoint,
		resource.Method,
		dbhelper.NullableString(resource.StatusCode),
	).Scan(&resourceID); err != nil {
		panic(fmt.Errorf("error finding resource: %w", err))
	}
	return resourceID
}

func findPropertyID(ctx context.Context, tx pgx.Tx, resourceID int64, path string) int64 {
	var propertyID int64
	if err := tx.QueryRow(ctx, findPropertyQuery, resourceID, path).Scan(&propertyID); err != nil {
		panic(fmt.Errorf("error finding property: %w", err))
	}
	return propertyID
}

func loadLatestContract(ctx context.Context, tx pgx.Tx, contractID int64, name, owner string) *model.Contract {
	contract := &model.Contract{Name: name, Owner: owner}

	rows, err := tx.Query(
		ctx,
		`SELECT id, direction, kind, provider, endpoint, method, status_code
		 FROM resources WHERE contract_id = $1`,
		contractID,
	)
	if err != nil {
		panic(fmt.Errorf("error loading resources: %w", err))
	}

	type resourceRow struct {
		id         int64
		direction  string
		kind       string
		provider   *string
		endpoint   string
		method     string
		statusCode *string
	}

	var resourceRows []resourceRow
	for rows.Next() {
		var rr resourceRow
		if err := rows.Scan(&rr.id, &rr.direction, &rr.kind, &rr.provider, &rr.endpoint, &rr.method, &rr.statusCode); err != nil {
			rows.Close()
			panic(fmt.Errorf("error scanning resource: %w", err))
		}
		resourceRows = append(resourceRows, rr)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("error iterating resources: %w", err))
	}

	for _, rr := range resourceRows {
		properties := loadLatestProperties(ctx, tx, rr.id)
		provider := ""
		if rr.provider != nil {
			provider = *rr.provider
		}
		statusCode := ""
		if rr.statusCode != nil {
			statusCode = *rr.statusCode
		}
		contract.AddResource(buildResource(rr.direction, rr.kind, provider, rr.endpoint, rr.method, statusCode, properties))
	}

	return contract
}

func loadLatestProperties(ctx context.Context, tx pgx.Tx, resourceID int64) map[string]model.Property {
	rows, err := tx.Query(
		ctx,
		`SELECT p.path, pv.type, pv.optional, pv.change
		 FROM properties p
		 JOIN LATERAL (
		     SELECT type, optional, change
		     FROM property_versions
		     WHERE property_id = p.id
		     ORDER BY contract_version_id DESC
		     LIMIT 1
		 ) pv ON true
		 WHERE p.resource_id = $1`,
		resourceID,
	)
	if err != nil {
		panic(fmt.Errorf("error loading properties: %w", err))
	}
	defer rows.Close()

	properties := map[string]model.Property{}
	for rows.Next() {
		var path, change string
		var propType *string
		var optional *bool
		if err := rows.Scan(&path, &propType, &optional, &change); err != nil {
			panic(fmt.Errorf("error scanning property: %w", err))
		}
		if change == string(model.ChangeRemoved) {
			continue
		}
		t := ""
		if propType != nil {
			t = *propType
		}
		o := false
		if optional != nil {
			o = *optional
		}
		properties[path] = model.NewProperty(path, t, o)
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("error iterating properties: %w", err))
	}
	return properties
}

func buildResource(direction, kind, provider, endpoint, method, statusCode string, properties map[string]model.Property) model.Resource {
	switch {
	case direction == string(model.Provides) && kind == string(model.RestResponse):
		return model.NewProvidedRestResponse(endpoint, method, statusCode, properties)
	case direction == string(model.Consumes) && kind == string(model.RestResponse):
		return model.NewConsumedRestResponse(provider, endpoint, method, statusCode, properties)
	case direction == string(model.Provides) && kind == string(model.RestRequest):
		return model.NewProvidedRestRequest(endpoint, method, properties)
	case direction == string(model.Consumes) && kind == string(model.RestRequest):
		return model.NewConsumedRestRequest(provider, endpoint, method, properties)
	}
	return model.Resource{
		Direction:  model.Direction(direction),
		Kind:       model.ResourceKind(kind),
		Provider:   provider,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		Properties: properties,
	}
}
