package model

type Contract struct {
	Schemas   map[string]Schema `json:"schemas,omitzero"`
	Resources []Resource        `json:"resources,omitzero"`
}
