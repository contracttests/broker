package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
	"github.com/contracttests/broker/internal/repository"
)

func ValidateContract(contract model.Contract) bool {
	hasError := false

	for _, resource := range contract.Resources {
		if resource.IsConsumer() {
			providerResource := repository.GetResource(resource.ProviderUuid)
			if providerResource.IsZero() {
				fmt.Println("Provider rest resource not found")
			}

			diff := SchemaDiff(resource.Schema, providerResource.Schema)

			if !diff.HasProperty() {
				continue
			}

			PrintInvalidConsumerResource(resource, providerResource, diff)
			hasError = true
		}

		if resource.IsProvider() {
			consumerResources := repository.GetConsumerResources(resource.ProviderUuid)
			if len(consumerResources) == 0 {
				continue
			}

			for _, consumerRestResource := range consumerResources {
				diff := SchemaDiff(consumerRestResource.Schema, resource.Schema)

				if !diff.HasProperty() {
					continue
				}

				PrintInvalidProviderResource(resource, consumerRestResource, diff)
				hasError = true
			}
		}
	}

	return hasError
}
