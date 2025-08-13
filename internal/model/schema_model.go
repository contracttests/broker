package model

type Property struct {
	Path string `json:"path,omitzero"`
	Type string `json:"type,omitzero"`
}

func NewProperty(path string, propertyType string) Property {
	return Property{
		Path: path,
		Type: propertyType,
	}
}

type Schema struct {
	Hash       string              `json:"hash,omitzero"`
	Properties map[string]Property `json:"properties,omitzero"`
}

func (s *Schema) IsZero() bool {
	return s.Hash == "" && len(s.Properties) == 0
}

func (s *Schema) HasProperty() bool {
	return len(s.Properties) > 0
}
