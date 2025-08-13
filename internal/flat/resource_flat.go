package flat

import (
	"strings"

	"github.com/contracttests/broker/internal/model"
)

type FlatResource struct {
	FullPath   string
	SchemaName string
}

func NewRestResource(
	resource FlatResource,
) model.RestResource {
	parts := strings.Split(resource.FullPath, ";")

	if strings.Contains(resource.FullPath, "consumes") {
		if strings.Contains(resource.FullPath, "requestBody") {
			consumerName, providerName, endpoint, method := parts[0], parts[2], parts[4], parts[5]

			return model.NewConsumerRestRequestBody(consumerName, providerName, endpoint, method)
		}

		consumerName, providerName, endpoint, method, statusCode := parts[0], parts[2], parts[4], parts[5], parts[7]

		return model.NewConsumerRestResponse(consumerName, providerName, endpoint, method, statusCode)
	}

	if strings.Contains(resource.FullPath, "requestBody") {
		providerName, endpoint, method := parts[0], parts[3], parts[4]

		return model.NewProviderRestRequestBody(providerName, endpoint, method)
	}

	providerName, endpoint, method, statusCode := parts[0], parts[3], parts[4], parts[6]

	return model.NewProviderRestResponse(providerName, endpoint, method, statusCode)
}
