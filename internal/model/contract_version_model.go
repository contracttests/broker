package model

import "github.com/google/uuid"

type ContractVersion struct {
	ID         int64
	UUID       uuid.UUID
	ContractID int64
	Version    int
	Checksum   string
	RawPayload string
}

func NewContractVersion(c *Contract) *ContractVersion {
	return &ContractVersion{
		ContractID: c.ID,
		Checksum:   c.Checksum(),
		RawPayload: c.RawPayload,
	}
}
