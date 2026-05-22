package wiredb

import (
	"github.com/contracttests/broker/server/internal/model"
	"github.com/google/uuid"
)

type InsertPropertyRow struct {
	ID         int64
	UUID       uuid.UUID
	ResourceID int64
	Path       string
}

func NewInsertPropertyRow(r model.Resource, p model.Property) *InsertPropertyRow {
	return &InsertPropertyRow{
		UUID:       uuid.New(),
		ResourceID: r.ID,
		Path:       p.Path,
	}
}
