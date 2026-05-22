package wiredb

import (
	"database/sql"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/google/uuid"
)

type InsertResourceRow struct {
	ID           int64
	UUID         uuid.UUID
	ContractID   int64
	Direction    string
	Kind         string
	Provider     sql.NullString
	Endpoint     string
	Method       string
	StatusCode   sql.NullString
	ProviderHash string
	ConsumerHash sql.NullString
}

func NewInsertResourceRow(c *model.Contract, r model.Resource) *InsertResourceRow {
	return &InsertResourceRow{
		UUID:         uuid.New(),
		ContractID:   c.ID,
		Direction:    string(r.Direction),
		Kind:         string(r.Kind),
		Provider:     sql.NullString{String: r.Provider, Valid: r.Provider != ""},
		Endpoint:     r.Endpoint,
		Method:       r.Method,
		StatusCode:   sql.NullString{String: r.StatusCode, Valid: r.StatusCode != ""},
		ProviderHash: r.ProviderHash(),
		ConsumerHash: sql.NullString{String: r.ConsumerHash(), Valid: r.ConsumerHash() != ""},
	}
}
