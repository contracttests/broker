package model

type Contract struct {
	Schemas       map[string]Schema `json:"schemas,omitzero"`
	RestResources []RestResource    `json:"restResources,omitzero"`
}
