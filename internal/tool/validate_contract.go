package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
	"github.com/contracttests/broker/internal/repository"
)

func ValidateContract(contract model.Contract) bool {
	hasError := false

	for _, restResource := range contract.Resources {
		if restResource.IsConsumer() {
			providerResource := repository.GetResource(restResource.ProviderUuid)
			if providerResource.IsZero() {
				fmt.Println("Provider rest resource not found")
			}

			leftSchema := contract.Schemas[restResource.SchemaUuid]
			rightSchema := repository.GetSchema(providerResource.SchemaUuid)

			diff := SchemaDiff(leftSchema, rightSchema)

			if !diff.HasProperty() {
				continue
			}

			PrintInvalidConsumerResource(restResource, providerResource, diff)
			hasError = true
		}

		if restResource.IsProvider() {
			consumerResources := repository.GetConsumerResources(restResource.ProviderUuid)
			if len(consumerResources) == 0 {
				continue
			}

			for _, consumerRestResource := range consumerResources {
				leftSchema := repository.GetSchema(consumerRestResource.SchemaUuid)
				rightSchema := contract.Schemas[restResource.SchemaUuid]

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
