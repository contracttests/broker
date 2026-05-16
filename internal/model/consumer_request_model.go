package model

type ConsumerRequest struct {
	Owner        string       `json:"owner"`
	Provider     string       `json:"provider"`
	Schema       Schema       `json:"schema"`
	ResourceType ResourceType `json:"resourceType"`
	RestRequest  RestRequest  `json:"restRequest,omitzero"`
}
