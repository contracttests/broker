package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
	"github.com/contracttests/broker/internal/repository"
)

func ValidateContract(contract model.Contract) bool {
	hasError := false

	for _, restResource := range contract.RestResources {
		if restResource.IsConsumer() {
			providerRestResource := repository.GetRestResource(restResource.ProviderHash)
			if providerRestResource.IsZero() {
				fmt.Println("Provider rest resource not found")
			}

			leftSchema := contract.Schemas[restResource.UniqueHash]
			rightSchema := repository.GetSchema(providerRestResource.UniqueHash)

			diff := SchemaDiff(leftSchema, rightSchema)

			if !diff.HasProperty() {
				continue
			}

			PrintInvalidConsumerResource(restResource, providerRestResource, diff)
			hasError = true
		}

		if restResource.IsProvider() {
			consumerRestResouces := repository.GetConsumerRestResources(restResource.ProviderHash)
			if len(consumerRestResouces) == 0 {
				continue
			}

			for _, consumerRestResource := range consumerRestResouces {
				leftSchema := repository.GetSchema(consumerRestResource.UniqueHash)
				rightSchema := contract.Schemas[restResource.UniqueHash]

				diff := SchemaDiff(leftSchema, rightSchema)

				if !diff.HasProperty() {
					continue
				}

				PrintInvalidProviderResource(restResource, consumerRestResource, diff)
				hasError = true
			}
		}
	}

	return hasError
}
