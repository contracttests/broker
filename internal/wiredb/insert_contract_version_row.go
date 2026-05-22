package wiredb

import (
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/google/uuid"
)

type InsertContractVersionRow struct {
	ID         int64
	UUID       uuid.UUID
	ContractID int64
	Version    int
	Checksum   string
	RawPayload string
}

func NewInsertContractVersionRow(cv *model.ContractVersion) *InsertContractVersionRow {
	return &InsertContractVersionRow{
		UUID:       uuid.New(),
		ContractID: cv.ContractID,
		Checksum:   cv.Checksum,
		RawPayload: cv.RawPayload,
	}
}
