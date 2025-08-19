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
	Properties map[string]Property `json:"properties,omitzero"`
}

func NewSchema() Schema {
	return Schema{
		Properties: make(map[string]Property),
	}
}

func (s *Schema) IsZero() bool {
	return len(s.Properties) == 0
}

func (s *Schema) HasProperty() bool {
	return len(s.Properties) > 0
}

type PropertyChange struct {
	ChangeType string   `json:"changeType,omitzero"`
	Property   Property `json:"property,omitzero"`
}

type SchemaChange struct {
	ChangeType string           `json:"changeType,omitzero"`
	Properties []PropertyChange `json:"properties,omitzero"`
}
