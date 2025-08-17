package repository

import (
	"github.com/contracttests/broker/internal/model"
)

var resourcesMap = make(map[string]model.Resource)
var consumerResouces = make(map[string][]model.Resource)

func SaveResource(resource model.Resource) {
	if resource.IsProvider() {
		resourcesMap[resource.Provider.Uuid] = resource
	} else {
		resourcesMap[resource.Consumer.Uuid] = resource
		consumerResouces[resource.Consumer.ProviderUuid] = append(consumerResouces[resource.Consumer.ProviderUuid], resource)
	}

}

func GetResource(uuid string) model.Resource {
	var resource model.Resource
	if resource, ok := resourcesMap[uuid]; ok {
		return resource
	}

	return resource
}

func GetConsumerResources(uuid string) []model.Resource {
	return consumerResouces[uuid]
}
