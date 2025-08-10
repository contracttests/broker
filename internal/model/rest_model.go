package model

import (
	"strings"

	"github.com/contracttests/broker/internal/flat"
)

type RestResource struct {
	UniqueHash   string `json:"uniqueHash,omitzero"`
	ProviderHash string `json:"providerHash,omitzero"`
	ConsumerName string `json:"consumerName,omitzero"`
	ProviderName string `json:"providerName,omitzero"`
	Endpoint     string `json:"endpoint,omitzero"`
	Method       string `json:"method,omitzero"`
	StatusCode   string `json:"statusCode,omitzero"`
	Direction    string `json:"direction,omitzero"`
	Type         string `json:"type,omitzero"`
}

func (rr *RestResource) IsZero() bool {
	return rr.UniqueHash == "" &&
		rr.ProviderHash == "" &&
		rr.ConsumerName == "" &&
		rr.ProviderName == "" &&
		rr.Endpoint == "" &&
		rr.Method == "" &&
		rr.StatusCode == "" &&
		rr.Direction == "" &&
		rr.Type == ""
}

func (rr *RestResource) IsProvider() bool {
	return rr.Type == "provider"
}

func (rr *RestResource) IsConsumer() bool {
	return rr.Type == "consumer"
}

func (rr *RestResource) IsRequestBody() bool {
	return rr.Direction == "requestBody"
}

func (rr *RestResource) IsResponse() bool {
	return rr.Direction == "response"
}

func NewRestResource(
	resource flat.FlatResource,
) RestResource {
	parts := strings.Split(resource.FullPath, ";")

	if strings.Contains(resource.FullPath, "consumes") {
		if strings.Contains(resource.FullPath, "requestBody") {
			consumerName, providerName, endpoint, method := parts[0], parts[2], parts[4], parts[5]

			return NewConsumerRestRequestBody(consumerName, providerName, endpoint, method)
		}

		consumerName, providerName, endpoint, method, statusCode := parts[0], parts[2], parts[4], parts[5], parts[7]

		return NewConsumerRestResponse(consumerName, providerName, endpoint, method, statusCode)
	}
	if strings.Contains(resource.FullPath, "requestBody") {
		providerName, endpoint, method := parts[0], parts[3], parts[4]

		return NewProviderRestRequestBody(providerName, endpoint, method)
	}

	providerName, endpoint, method, statusCode := parts[0], parts[3], parts[4], parts[6]

	return NewProviderRestResponse(providerName, endpoint, method, statusCode)
}

func NewConsumerRestRequestBody(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
) RestResource {
	uniqueHash := HashFromStrings(consumerName, "consumes", providerName, endpoint, method, "requestBody")
	providerHash := HashFromStrings(providerName, "provides", endpoint, method, "requestBody")

	return RestResource{
		UniqueHash:   uniqueHash,
		ProviderHash: providerHash,
		ConsumerName: consumerName,
		ProviderName: providerName,
		Endpoint:     endpoint,
		Method:       method,
		Direction:    "requestBody",
		Type:         "consumer",
	}
}

func NewConsumerRestResponse(
	consumerName string,
	providerName string,
	endpoint string,
	method string,
	statusCode string,
) RestResource {
	uniqueHash := HashFromStrings(consumerName, "consumes", providerName, endpoint, method, statusCode, "response")
	providerHash := HashFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	return RestResource{
		UniqueHash:   uniqueHash,
		ProviderHash: providerHash,
		ConsumerName: consumerName,
		ProviderName: providerName,
		Endpoint:     endpoint,
		Method:       method,
		StatusCode:   statusCode,
		Direction:    "response",
		Type:         "consumer",
	}
}

func NewProviderRestRequestBody(
	providerName string,
	endpoint string,
	method string,
) RestResource {
	uniqueHash := HashFromStrings(providerName, "provides", endpoint, method, "requestBody")

	return RestResource{
		UniqueHash:   uniqueHash,
		ProviderHash: uniqueHash,
		ProviderName: providerName,
		Endpoint:     endpoint,
		Method:       method,
		Direction:    "requestBody",
		Type:         "provider",
	}
}

func NewProviderRestResponse(
	providerName string,
	endpoint string,
	method string,
	statusCode string,
) RestResource {
	uniqueHash := HashFromStrings(providerName, "provides", endpoint, method, statusCode, "response")

	return RestResource{
		UniqueHash:   uniqueHash,
		ProviderHash: uniqueHash,
		ProviderName: providerName,
		Endpoint:     endpoint,
		Method:       method,
		StatusCode:   statusCode,
		Direction:    "response",
		Type:         "provider",
	}
}
