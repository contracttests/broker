package model

type ProviderRequest struct {
	Owner        string       `json:"owner"`
	Schema       Schema       `json:"schema"`
	ResourceType ResourceType `json:"resourceType"`
	RestRequest  RestRequest  `json:"restRequest,omitzero"`
}
