package flat

import (
	"strings"

	"github.com/contracttests/broker/internal/model"
)

type FlatContract struct {
	Resources []FlatResource
	Schemas   FlatSchemas
}

func (flatContract FlatContract) ToModelContract() model.Contract {
	contract := model.Contract{
		RestResources: []model.RestResource{},
		Schemas:       make(map[string]model.Schema),
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
