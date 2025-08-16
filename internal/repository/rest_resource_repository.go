package repository

import "github.com/contracttests/broker/internal/model"

var restResourcesMap = make(map[string]model.Resource)
var consumerRestResouce = make(map[string][]model.Resource)

func SaveResource(resource model.Resource) {
	if resource.ConsumerUuid == "" {
		restResourcesMap[resource.ProviderUuid] = resource
	} else {
		restResourcesMap[resource.ConsumerUuid] = resource
	}

	if resource.IsConsumer() {
		consumerRestResouce[resource.ProviderUuid] = append(consumerRestResouce[resource.ProviderUuid], resource)
	}
}

func GetResource(uuid string) model.Resource {
	var restResource model.Resource
	if restResource, ok := restResourcesMap[uuid]; ok {
		return restResource
	}

	return restResource
}

func GetConsumerResources(providerHash string) []model.Resource {
	return consumerRestResouce[providerHash]
}
