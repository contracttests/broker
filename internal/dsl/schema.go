package dsl

import "fmt"

type Properties map[string]Schema

type Schemas map[string]Schema

type Schema struct {
	Type       string     `json:"type,omitzero"`
	Properties Properties `json:"properties,omitzero"`
	Items      *Schema    `json:"items,omitzero"`
	Ref        string     `json:"$ref,omitzero"`
}

func (s *Schema) IsObject() bool {
	if s.Type != "" {
		return s.Type == "object"
	}

	return s.Properties != nil
}

func (s *Schema) IsArray() bool {
	if s.Type != "" {
		return s.Type == "array"
	}

	return s.Items != nil
}

func (s *Schema) IsPrimitive() bool {
	if s.Type != "" {
		return s.Type == "string" || s.Type == "integer" || s.Type == "float" || s.Type == "number" || s.Type == "boolean"
	}

	return false
}

func (s *Schema) IsRef() bool {
	if s.Type != "" || s.Properties != nil || s.Items != nil {
		return false
	}

	return s.Ref != ""
}

type ComparableSchemaProperty struct {
	Path string
	Type string
}

type ComparableSchemaProperties []ComparableSchemaProperty

type ComparableSchema struct {
	Name       string
	Properties ComparableSchemaProperties
}

type ComparableSchemas []ComparableSchema

func NewComparableSchemaProperty(path string, pathType string) ComparableSchemaProperty {
	return ComparableSchemaProperty{
		Path: path,
		Type: pathType,
	}
}

func NewComparableSchemas(schemas Schemas) ComparableSchemas {
	comparableSchemas := ComparableSchemas{}

	for name, schema := range schemas {
		comparableSchemas = append(comparableSchemas, ComparableSchema{
			Name:       name,
			Properties: newComparableSchemas([]ComparableSchemaProperty{}, "root", schema),
		})
	}

	return comparableSchemas
}

func newComparableSchemas(properties []ComparableSchemaProperty, path string, unknown any) []ComparableSchemaProperty {
	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			properties = append(properties, NewComparableSchemaProperty(path, "object"))

			for name, schema := range unknown.Properties {
				pathName := fmt.Sprintf("%s.%s", path, name)
				properties = newComparableSchemas(properties, pathName, schema)
			}

			return properties
		}

		if unknown.IsArray() {
			properties = append(properties, NewComparableSchemaProperty(path, "array"))
			pathName := fmt.Sprintf("%s[]", path)
			properties = newComparableSchemas(properties, pathName, unknown.Items)
			return properties
		}

		if unknown.IsPrimitive() {
			properties = append(properties, NewComparableSchemaProperty(path, unknown.Type))
			return properties
		}

		if unknown.IsRef() {
			properties = append(properties, NewComparableSchemaProperty(path, unknown.Ref))
			return properties
		}
	case *Schema:
		return newComparableSchemas(properties, path, *unknown)

	default:
		panic("unknown schema type")
	}

	return properties
}
