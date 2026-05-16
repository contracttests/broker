package model

type Contract struct {
	ConsumerRequests  []ConsumerRequest  `json:"consumerRequests,omitzero"`
	ConsumerResponses []ConsumerResponse `json:"consumerResponses,omitzero"`
	ProviderRequests  []ProviderRequest  `json:"providerRequests,omitzero"`
	ProviderResponses []ProviderResponse `json:"providerResponses,omitzero"`
}
