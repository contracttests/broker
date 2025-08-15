package repository

import "github.com/contracttests/broker/internal/model"

var restResourcesMap = make(map[string]model.Resource)
var consumerRestResouce = make(map[string][]model.Resource)

func SaveRestResource(resource model.Resource) {
	restResourcesMap[resource.Uuid] = resource

	if resource.IsConsumer() {
		consumerRestResouce[resource.ProviderUuid] = append(consumerRestResouce[resource.ProviderUuid], resource)
	}
}

func GetResource(hash string) model.Resource {
	var restResource model.Resource
	if restResource, ok := restResourcesMap[hash]; ok {
		return restResource
	}

	return restResource
}

func GetConsumerResources(providerHash string) []model.Resource {
	return consumerRestResouce[providerHash]
}
