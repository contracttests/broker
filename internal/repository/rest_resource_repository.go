package repository

import "github.com/contracttests/broker/internal/model"

var restResourcesMap = make(map[string]model.RestResource)
var consumerRestResouce = make(map[string][]model.RestResource)

func SaveRestResource(restResource model.RestResource) {
	restResourcesMap[restResource.UniqueHash] = restResource

	if restResource.IsConsumer() {
		consumerRestResouce[restResource.ProviderHash] = append(consumerRestResouce[restResource.ProviderHash], restResource)
	}
}

func GetRestResource(hash string) model.RestResource {
	var restResource model.RestResource
	if restResource, ok := restResourcesMap[hash]; ok {
		return restResource
	}

	return restResource
}

func GetConsumerRestResources(providerHash string) []model.RestResource {
	return consumerRestResouce[providerHash]
}
