package model

type Property struct {
	Path     string `json:"path,omitzero"`
	Type     string `json:"type,omitzero"`
	Optional bool   `json:"optional"`
}

func (p *Property) IsSame(pp *Property) bool {
	return p.Path == pp.Path &&
		p.Type == pp.Type &&
		p.Optional == pp.Optional
}

func NewProperty(
	propertyPath string,
	propertyType string,
	optional bool,
) Property {
	return Property{
		Path:     propertyPath,
		Type:     propertyType,
		Optional: optional,
	}
}
