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
	Consumer     ConsumerResource `json:"consumer,omitzero"`
	Provider     ProviderResource `json:"provider,omitzero"`
	Schema       Schema           `json:"schema,omitzero"`
	Direction    string           `json:"direction,omitzero"`
	RestResource RestResource     `json:"restResource,omitzero"`
}

func (r *Resource) IsConsumer() bool {
	return r.Consumer.Name != ""
}

func (r *Resource) IsProvider() bool {
	return r.Provider.Name != ""
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
	return rr.Consumer == ConsumerResource{} &&
		rr.Provider == ProviderResource{} &&
		rr.Direction == "" &&
		len(rr.Schema.Properties) == 0
}

func (rr *Resource) IsRequestBody() bool {
	return rr.Direction == "request"
}

func NewResource(fullPath string, schema Schema) Resource {
	parts := strings.Split(fullPath, ";")

	if strings.Contains(fullPath, "consumes") {
		if strings.Contains(fullPath, "request") {
			consumerName, providerName, endpoint, method := parts[0], parts[2], parts[4], parts[5]

			return NewConsumerRequestBody(consumerName, providerName, endpoint, method, schema)
		}

		consumerName, providerName, endpoint, method, statusCode := parts[0], parts[2], parts[4], parts[5], parts[7]

		return NewConsumerResponse(consumerName, providerName, endpoint, method, statusCode, schema)
	}

	if strings.Contains(fullPath, "request") {
		providerName, endpoint, method := parts[0], parts[3], parts[4]

		return NewProviderRequestBody(providerName, endpoint, method, schema)
	}

	providerName, endpoint, method, statusCode := parts[0], parts[3], parts[4], parts[6]

	return NewProviderResponse(providerName, endpoint, method, statusCode, schema)
}

func NewConsumerRequestBody(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method)
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method)

	return Resource{
		Consumer: ConsumerResource{
			Uuid:         uuid,
			Name:         consumerName,
			ProviderUuid: providerUuid,
		},
		Schema: schema,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
	}
}

func NewConsumerResponse(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	statusCode string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method, statusCode)
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode)

	return Resource{
		Consumer: ConsumerResource{
			Uuid:         uuid,
			Name:         consumerName,
			ProviderUuid: providerUuid,
		},
		Schema: schema,
		RestResource: RestResource{
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: statusCode,
		},
		Direction: "response",
	}
}

func NewProviderRequestBody(
	providerName string,
	endpoint string,
	method string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method)

	return Resource{
		Provider: ProviderResource{
			Uuid: uuid,
			Name: providerName,
		},
		Schema: schema,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
	}
}

func NewProviderResponse(
	providerName string,
	endpoint string,
	method string,
	statusCode string,
	schema Schema,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode)

	return Resource{
		Provider: ProviderResource{
			Uuid: uuid,
			Name: providerName,
		},
		Schema: schema,
		RestResource: RestResource{
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: statusCode,
		},
		Direction: "response",
	}
}
