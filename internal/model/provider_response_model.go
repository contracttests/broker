package model

type ProviderResponse struct {
	Owner        string       `json:"owner"`
	Schema       Schema       `json:"schema"`
	ResourceType ResourceType `json:"resourceType"`
	RestResponse RestResponse `json:"restResponse,omitzero"`
}
