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
			providerResource := repository.GetResource(resource.Consumer.ProviderUuid)
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
			consumerResources := repository.GetConsumerResources(resource.Provider.Uuid)
			if len(consumerResources) == 0 {
				continue
			}

			for _, consumerResource := range consumerResources {
				diff := SchemaDiff(consumerResource.Schema, resource.Schema)

				if !diff.HasProperty() {
					continue
				}

				PrintInvalidProviderResource(resource, consumerResource, diff)
				hasError = true
			}
		}
	}

	return hasError
}
