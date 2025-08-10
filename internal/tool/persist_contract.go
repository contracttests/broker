package tool

import (
	"github.com/contracttests/broker/internal/model"
	"github.com/contracttests/broker/internal/repository"
)

func SaveContract(contract model.Contract) {
	for _, restResource := range contract.RestResources {
		repository.SaveRestResource(restResource)
	}

	for _, schema := range contract.Schemas {
		repository.SaveSchema(schema)
	}
}
