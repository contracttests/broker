package model

import (
	"strings"
)

type RestResource struct {
	Endpoint   string `json:"endpoint,omitzero"`
	Method     string `json:"method,omitzero"`
	StatusCode string `json:"statusCode,omitzero"`
}

type Resource struct {
	ConsumerUuid string           `json:"consumerUuid,omitzero"`
	ProviderUuid string           `json:"providerUuid,omitzero"`
	ConsumerName string           `json:"consumerName,omitzero"`
	ProviderName string           `json:"providerName,omitzero"`
	Consumer     ConsumerResource `json:"consumer,omitzero"`
	Provider     ProviderResource `json:"provider,omitzero"`
	Schema       Schema           `json:"schema,omitzero"`
	Direction    string           `json:"direction,omitzero"`
	Type         string           `json:"type,omitzero"`
	RestResource RestResource     `json:"restResource,omitzero"`
}

type ConsumerResource struct {
	Uuid         string `json:"uuid,omitzero"`
	Name         string `json:"name,omitzero"`
	ProviderUuid string `json:"providerUuid,omitzero"`
}

type ProviderResource struct {
	Uuid string `json:"uuid,omitzero"`
	Name string `json:"name,omitzero"`
}

func (rr *Resource) IsZero() bool {
	return rr.ConsumerUuid == "" &&
		rr.ProviderUuid == "" &&
		rr.ConsumerName == "" &&
		rr.ProviderName == "" &&
		rr.Direction == "" &&
		rr.Type == ""
}

func (rr *Resource) IsProvider() bool {
	return rr.Type == "provider"
}

func (rr *Resource) IsConsumer() bool {
	return rr.Type == "consumer"
}

func (rr *Resource) IsRequestBody() bool {
	return rr.Direction == "request"
}

func NewResource(fullPath string, schema Schema) Resource {
	parts := strings.Split(fullPath, ";")

	if strings.Contains(fullPath, "consumes") {
		if strings.Contains(fullPath, "request") {
			consumerName, providerName, endpoint, method := parts[0], parts[2], parts[4], parts[5]

			return NewConsumerRestRequestBody(consumerName, providerName, endpoint, method, schema)
		}

		consumerName, providerName, endpoint, method, statusCode := parts[0], parts[2], parts[4], parts[5], parts[7]

		return NewConsumerRestResponse(consumerName, providerName, endpoint, method, statusCode, schema)
	}

	if strings.Contains(fullPath, "request") {
		providerName, endpoint, method := parts[0], parts[3], parts[4]

		return NewProviderRestRequestBody(providerName, endpoint, method, schema)
	}

	providerName, endpoint, method, statusCode := parts[0], parts[3], parts[4], parts[6]

	return NewProviderRestResponse(providerName, endpoint, method, statusCode, schema)
}

func NewConsumerRestRequestBody(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method, "request")
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method, "request")

	consumer := ConsumerResource{
		Uuid:         uuid,
		Name:         consumerName,
		ProviderUuid: providerUuid,
	}

	return Resource{
		ConsumerUuid: uuid,
		ProviderUuid: providerUuid,
		ConsumerName: consumerName,
		ProviderName: providerName,
		Consumer:     consumer,
		Schema:       schema,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
		Type:      "consumer",
	}
}

func NewConsumerRestResponse(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	statusCode string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method, statusCode, "response")
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	consumer := ConsumerResource{
		Uuid:         uuid,
		Name:         consumerName,
		ProviderUuid: providerUuid,
	}

	return Resource{
		ConsumerUuid: uuid,
		ProviderUuid: providerUuid,
		ConsumerName: consumerName,
		ProviderName: providerName,
		Consumer:     consumer,
		Schema:       schema,
		RestResource: RestResource{
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: statusCode,
		},
		Direction: "response",
		Type:      "consumer",
	}
}

func NewProviderRestRequestBody(
	providerName string,
	endpoint string,
	method string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method, "request")

	provider := ProviderResource{
		Uuid: uuid,
		Name: providerName,
	}

	return Resource{
		ProviderUuid: uuid,
		ProviderName: providerName,
		Provider:     provider,
		Schema:       schema,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
		Type:      "provider",
	}
}

func NewProviderRestResponse(
	providerName string,
	endpoint string,
	method string,
	statusCode string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	provider := ProviderResource{
		Uuid: uuid,
		Name: providerName,
	}

	return Resource{
		ProviderUuid: uuid,
		ProviderName: providerName,
		Provider:     provider,
		Schema:       schema,
		RestResource: RestResource{
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: statusCode,
		},
		Direction: "response",
		Type:      "provider",
	}
}
