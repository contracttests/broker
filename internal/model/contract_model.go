package model

type Contract struct {
	Resources []Resource        `json:"resources,omitzero"`
	Schemas   map[string]Schema `json:"schemas,omitzero"`
}
