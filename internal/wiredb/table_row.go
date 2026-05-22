package wiredb

import (
	"database/sql"
	"time"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/google/uuid"
)

type TableRow struct {
	ContractID        int64
	ContractUUID      uuid.UUID
	ContractName      string
	ContractOwner     string
	ContractCreatedAt time.Time

	ResourceID           int64
	ResourceUUID         uuid.UUID
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

func (c *TableRow) ToContractModel() *model.Contract {
	return &model.Contract{
		ID:        c.ContractID,
		UUID:      c.ContractUUID,
		Name:      c.ContractName,
		Owner:     c.ContractOwner,
		Resources: make(map[string]model.Resource),
	}
}

func (c *TableRow) ToResourceModel() model.Resource {
	return model.Resource{
		ID:         c.ResourceID,
		Direction:  model.Direction(c.ResourceDirection),
		Kind:       model.ResourceKind(c.ResourceKind),
		Provider:   c.ResourceProvider.String,
		Endpoint:   c.ResourceEndpoint,
		Method:     c.ResourceMethod,
		StatusCode: c.ResourceStatusCode.String,
		Properties: make(map[string]model.Property),
		ContractInfo: &model.ContractInfo{
			Name:  c.ContractName,
			Owner: c.ContractOwner,
		},
	}
}

func (c *TableRow) ToPropertyModel() model.Property {
	return model.Property{
		ID:       c.PropertyID,
		Path:     c.PropertyPath,
		Type:     c.PropertyVersionType.String,
		Optional: c.PropertyVersionOptional.Bool,
	}
}
