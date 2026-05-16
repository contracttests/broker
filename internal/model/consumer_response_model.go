package model

type ConsumerResponse struct {
	Owner        string       `json:"owner"`
	Provider     string       `json:"provider"`
	Schema       Schema       `json:"schema"`
	ResourceType ResourceType `json:"resourceType"`
	RestResponse RestResponse `json:"restResponse,omitzero"`
}
