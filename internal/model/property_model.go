package model

import (
	"strconv"
	"strings"
)

type Property struct {
	ID       int64
	Path     string
	Type     string
	Optional bool
}

func (p *Property) IsSame(other *Property) bool {
	return p.Path == other.Path &&
		p.Type == other.Type &&
		p.Optional == other.Optional
}

func (p *Property) CanonicalKey() string {
	return strings.Join([]string{
		p.Path,
		p.Type,
		strconv.FormatBool(p.Optional),
	}, ";;")
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
