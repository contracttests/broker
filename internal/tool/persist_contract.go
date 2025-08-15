package tool

import (
	"github.com/contracttests/broker/internal/model"
	"github.com/contracttests/broker/internal/repository"
)

func SaveContract(contract model.Contract) {
	for _, resource := range contract.Resources {
		repository.SaveRestResource(resource)
	}

	for _, schema := range contract.Schemas {
		repository.SaveSchema(schema)
	}
}
