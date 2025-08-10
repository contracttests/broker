package model

import (
	"strings"

	"github.com/contracttests/broker/internal/flat"
)

type Contract struct {
	Schemas       map[string]Schema `json:"schemas,omitzero"`
	RestResources []RestResource    `json:"restResources,omitzero"`
}

func NewContract(
	flatContract flat.FlatContract,
) Contract {
	contract := Contract{
		RestResources: []RestResource{},
		Schemas:       make(map[string]Schema),
	}

	for _, flatResource := range flatContract.Resources {
		if strings.Contains(flatResource.FullPath, "rest") {
			resource := NewRestResource(flatResource)
			contract.RestResources = append(contract.RestResources, resource)

			flatSchema := flatContract.Schemas[flatResource.SchemaName]

			schema := NewSchema(resource.UniqueHash, flatSchema)
			contract.Schemas[resource.UniqueHash] = schema
		}
	}

	return contract
}
