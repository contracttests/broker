package repository

import (
	"database/sql"
	"time"

	"github.com/contracttesting/broker/server/internal/model"
)

type tableRow struct {
	ParticipantID   int64
	ParticipantName string

	ContractID          int64
	ContractVersion     int
	ContractRawContract string
	ContractCreatedAt   time.Time

	ResourceID           int64
	ResourceDirection    string
	ResourceKind         string
	ResourceProvider     sql.NullString
	ResourceEndpoint     string
	ResourceMethod       string
	ResourceStatusCode   sql.NullString
	ResourceProviderHash string
	ResourceConsumerHash sql.NullString
	ResourceCreatedAt    time.Time

	PropertyID   int64
	PropertyPath string

	PropertyVersionType     sql.NullString
	PropertyVersionOptional sql.NullBool
	PropertyVersionChange   string
}

func (c *tableRow) toContractModel() *model.Contract {
	return &model.Contract{
		ID:          c.ContractID,
		Version:     c.ContractVersion,
		RawContract: c.ContractRawContract,
		Resources:   make(map[string]model.Resource),
		Participant: &model.Participant{
			ID:   c.ParticipantID,
			Name: c.ParticipantName,
		},
	}
}

func (c *tableRow) toResourceModel() model.Resource {
	return model.Resource{
		ID:         c.ResourceID,
		Direction:  model.Direction(c.ResourceDirection),
		Kind:       model.ResourceKind(c.ResourceKind),
		Provider:   c.ResourceProvider.String,
		Endpoint:   c.ResourceEndpoint,
		Method:     c.ResourceMethod,
		StatusCode: c.ResourceStatusCode.String,
		Properties: make(map[string]model.Property),
		Participant: &model.Participant{
			ID:   c.ParticipantID,
			Name: c.ParticipantName,
		},
	}
}

func (c *tableRow) toPropertyModel() model.Property {
	return model.Property{
		ID:       c.PropertyID,
		Path:     c.PropertyPath,
		Type:     c.PropertyVersionType.String,
		Optional: c.PropertyVersionOptional.Bool,
	}
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

type insertPropertyVersionRow struct {
	PropertyID int64
	ContractID int64
	Type       sql.NullString
	Optional   sql.NullBool
	Change     string
}

func newInsertPropertyVersionRowAdded(c *model.Contract, p model.Property) *insertPropertyVersionRow {
	return newInsertPropertyVersionRow(c, p, model.ChangeAdded)
}

func newInsertPropertyVersionRowModified(c *model.Contract, p model.Property) *insertPropertyVersionRow {
	return newInsertPropertyVersionRow(c, p, model.ChangeModified)
}

func newInsertPropertyVersionRowRemoved(c *model.Contract, p model.Property) *insertPropertyVersionRow {
	return newInsertPropertyVersionRow(c, p, model.ChangeRemoved)
}

func newInsertPropertyVersionRow(c *model.Contract, p model.Property, change model.ChangeKind) *insertPropertyVersionRow {
	return &insertPropertyVersionRow{
		PropertyID: p.ID,
		ContractID: c.ID,
		Type:       nullString(p.Type),
		Optional:   sql.NullBool{Bool: p.Optional, Valid: true},
		Change:     string(change),
	}
}

type insertResourceVersionRow struct {
	ResourceID int64
	ContractID int64
	Change     string
}

func newInsertResourceVersionRowAdded(c *model.Contract, r model.Resource) *insertResourceVersionRow {
	return &insertResourceVersionRow{
		ResourceID: r.ID,
		ContractID: c.ID,
		Change:     string(model.ChangeAdded),
	}
}

func newInsertResourceVersionRowRemoved(c *model.Contract, r model.Resource) *insertResourceVersionRow {
	return &insertResourceVersionRow{
		ResourceID: r.ID,
		ContractID: c.ID,
		Change:     string(model.ChangeRemoved),
	}
}
