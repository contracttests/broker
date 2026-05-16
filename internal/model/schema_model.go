package model

type SchemaProperties map[string]Property

type Schema struct {
	Properties SchemaProperties `json:"properties,omitzero"`
}

func (s *Schema) AddProperty(property Property) {
	s.Properties[property.Path] = property
}

func NewSchema() Schema {
	return Schema{
		Properties: make(SchemaProperties),
	}
}
