package wiredb

import (
	"database/sql"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/google/uuid"
)

type InsertPropertyVersionRow struct {
	UUID              uuid.UUID
	PropertyID        int64
	ContractVersionID int64
	Type              sql.NullString
	Optional          sql.NullBool
	Change            string
}

func NewInsertPropertyVersionRowAdded(cv *model.ContractVersion, p model.Property) *InsertPropertyVersionRow {
	return newInsertPropertyVersionRow(cv, p, model.ChangeAdded)
}

func NewInsertPropertyVersionRowModified(cv *model.ContractVersion, p model.Property) *InsertPropertyVersionRow {
	return newInsertPropertyVersionRow(cv, p, model.ChangeModified)
}

func NewInsertPropertyVersionRowRemoved(cv *model.ContractVersion, p model.Property) *InsertPropertyVersionRow {
	return newInsertPropertyVersionRow(cv, p, model.ChangeRemoved)
}

func newInsertPropertyVersionRow(cv *model.ContractVersion, p model.Property, change model.ChangeKind) *InsertPropertyVersionRow {
	return &InsertPropertyVersionRow{
		UUID:              uuid.New(),
		PropertyID:        p.ID,
		ContractVersionID: cv.ID,
		Type:              sql.NullString{String: p.Type, Valid: p.Type != ""},
		Optional:          sql.NullBool{Bool: p.Optional, Valid: true},
		Change:            string(change),
	}
}
