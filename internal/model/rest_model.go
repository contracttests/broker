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
	Uuid         string       `json:"uuid,omitzero"`
	ProviderUuid string       `json:"providerUuid,omitzero"`
	ConsumerName string       `json:"consumerUuid,omitzero"`
	ProviderName string       `json:"providerName,omitzero"`
	SchemaUuid   string       `json:"schemaUuid,omitzero"`
	Schema       Schema       `json:"schema,omitzero"`
	Direction    string       `json:"direction,omitzero"`
	Type         string       `json:"type,omitzero"`
	RestResource RestResource `json:"restResource,omitzero"`
}

func (rr *Resource) IsZero() bool {
	return rr.Uuid == "" &&
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

func NewResource(fullPath string, schemaUuid string) Resource {
	parts := strings.Split(fullPath, ";")

	if strings.Contains(fullPath, "consumes") {
		if strings.Contains(fullPath, "request") {
			consumerName, providerName, endpoint, method := parts[0], parts[2], parts[4], parts[5]

			return NewConsumerRestRequestBody(schemaUuid, consumerName, providerName, endpoint, method)
		}

		consumerName, providerName, endpoint, method, statusCode := parts[0], parts[2], parts[4], parts[5], parts[7]

		return NewConsumerRestResponse(schemaUuid, consumerName, providerName, endpoint, method, statusCode)
	}

	if strings.Contains(fullPath, "request") {
		providerName, endpoint, method := parts[0], parts[3], parts[4]

		return NewProviderRestRequestBody(schemaUuid, providerName, endpoint, method)
	}

	providerName, endpoint, method, statusCode := parts[0], parts[3], parts[4], parts[6]

	return NewProviderRestResponse(schemaUuid, providerName, endpoint, method, statusCode)
}

func NewConsumerRestRequestBody(
	schemaUuid string,
	consumerName string,
	providerName string,
	endpoint string,
	method string,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method, "request")
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method, "request")

	return Resource{
		Uuid:         uuid,
		SchemaUuid:   schemaUuid,
		ProviderUuid: providerUuid,
		ConsumerName: consumerName,
		ProviderName: providerName,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
		Type:      "consumer",
	}
}

func NewConsumerRestResponse(
	schemaUuid string,
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	statusCode string,
) Resource {
	uuid := UuidFromStrings(consumerName, "consumes", providerName, endpoint, method, statusCode, "response")
	providerUuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	return Resource{
		Uuid:         uuid,
		ProviderUuid: providerUuid,
		SchemaUuid:   schemaUuid,
		ConsumerName: consumerName,
		ProviderName: providerName,
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
	schemaUuid string,
	providerName string,
	endpoint string,
	method string,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method, "request")

	return Resource{
		Uuid:         uuid,
		SchemaUuid:   schemaUuid,
		ProviderUuid: uuid,
		ProviderName: providerName,
		RestResource: RestResource{
			Endpoint: endpoint,
			Method:   method,
		},
		Direction: "request",
		Type:      "provider",
	}
}

func NewProviderRestResponse(
	schemaUuid string,
	providerName string,
	endpoint string,
	method string,
	statusCode string,
) Resource {
	uuid := UuidFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	return Resource{
		Uuid:         uuid,
		ProviderUuid: uuid,
		SchemaUuid:   schemaUuid,
		ProviderName: providerName,
		RestResource: RestResource{
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: statusCode,
		},
		Direction: "response",
		Type:      "provider",
	}
}
