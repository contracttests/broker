package wiredb

import (
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/google/uuid"
)

type InsertContractRow struct {
	ID    int64
	UUID  uuid.UUID
	Name  string
	Owner string
}

func NewInsertContractRow(c *model.Contract) *InsertContractRow {
	return &InsertContractRow{
		UUID:  uuid.New(),
		Name:  c.Name,
		Owner: c.Owner,
	}
}
