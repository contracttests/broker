package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/contracttests/broker/internal/model"
)

type Api struct {
	Name string `json:"name,omitzero"`
}

type Contract struct {
	Api              Api                 `json:"api,omitzero"`
	Provides         Provides            `json:"provides,omitzero"`
	ConsumesServices map[string]Consumes `json:"consumes,omitzero"`
	Schemas          map[string]Schema   `json:"schemas,omitzero"`
}

type ResourcePath string

func (f *ResourcePath) Append(parts ...string) ResourcePath {
	separator := ";"

	return ResourcePath(strings.Join([]string{f.String(), strings.Join(parts, separator)}, separator))
}

func (f *ResourcePath) String() string {
	return string(*f)
}

type PropertyPath string

func (f *PropertyPath) Append(parts ...string) PropertyPath {
	separator := "."

	return PropertyPath(strings.Join([]string{f.String(), strings.Join(parts, separator)}, separator))
}

func (f *PropertyPath) String() string {

	return string(*f)
}

func (c *Contract) ToContractModel() model.Contract {
	resources := Resources(*c)
	schemas := Schemas(*c)

	return model.Contract{
		Resources: resources,
		Schemas:   schemas,
	}
}

func newResourceFullPath(parts ...string) string {
	resourceParthSeparator := ";"

	return strings.Join(parts, resourceParthSeparator)
}

func Resources(contractDsl Contract) []model.Resource {
	return buildResources([]model.Resource{}, "", contractDsl)
}

func buildResources(flatResources []model.Resource, fullPath string, unknown any) []model.Resource {
	switch unknown := unknown.(type) {
	case Contract:
		fullPath = newFullPath(fullPath, unknown.Api.Name)

		for serviceName, consumes := range unknown.ConsumesServices {
			newFullPath := newResourceFullPath(fullPath, "consumes", serviceName)
			flatResources = buildResources(flatResources, newFullPath, consumes)
		}

		newFullPath := newResourceFullPath(fullPath, "provides")
		flatResources = buildResources(flatResources, newFullPath, unknown.Provides)

		return flatResources

	case Consumes:
		flatResources = buildResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case Provides:
		flatResources = buildResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case Message:
		for messageName, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, "message", messageName)
			flatResources = append(flatResources, model.NewResource(newFullPath, model.UuidFromStrings(schemaName)))
		}

		return flatResources

	case Rest:
		for endpoint, methods := range unknown {
			if methods.Get.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Get)
			}

			if methods.Post.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Post)
			}

			if methods.Put.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Put)
			}

			if methods.Delete.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Delete)
			}
		}

	case GetMethod:
		newFullPath := newResourceFullPath(fullPath, "get", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PostMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "post", "request")
			flatResources = append(flatResources, model.NewResource(newFullPath, model.UuidFromStrings(unknown.RequestBody)))
		}

		newFullPath := newResourceFullPath(fullPath, "post", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PutMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "put", "request")
			flatResources = append(flatResources, model.NewResource(newFullPath, model.UuidFromStrings(unknown.RequestBody)))
		}

		newFullPath := newResourceFullPath(fullPath, "put", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case DeleteMethod:
		newFullPath := newResourceFullPath(fullPath, "delete", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case Responses:
		for statusCode, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, strconv.Itoa(statusCode))
			flatResources = append(flatResources, model.NewResource(newFullPath, model.UuidFromStrings(schemaName)))
		}

		return flatResources
	}

	return flatResources
}

func newFullPath(parts ...string) string {
	sanitizedParts := []string{}
	for _, part := range parts {
		if part != "" {
			sanitizedParts = append(sanitizedParts, part)
		}
	}

	if len(sanitizedParts) == 1 {
		return sanitizedParts[0]
	}

	return strings.Join(sanitizedParts, ".")
}

func newArrayPropertyPath(part string) string {
	return part + "[]"
}

func Schemas(contractDsl Contract) map[string]model.Schema {
	schemas := make(map[string]model.Schema)

	for schemaName, schema := range contractDsl.Schemas {
		hash := model.UuidFromStrings(schemaName)

		schema := buildSchema(
			0,
			schemaName,
			contractDsl.Schemas,
			model.Schema{
				Hash:       hash,
				Properties: make(map[string]model.Property),
			},
			"root",
			schema,
		)

		schemas[hash] = schema
	}

	return schemas
}

func buildSchema(
	deep int,
	originalSchemaName string,
	schemas map[string]Schema,
	schema model.Schema,
	fullPath string,
	unknown any,
) model.Schema {
	if deep >= 10 {
		panic(fmt.Sprintf("Circular reference detected in the schema %s", originalSchemaName))
	}

	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			schema.Properties[fullPath] = model.Property{
				Path: fullPath,
				Type: "object",
			}

			for name, schemaProperties := range unknown.Properties {
				schema = buildSchema(
					deep+1,
					originalSchemaName,
					schemas,
					schema,
					newFullPath(fullPath, name),
					schemaProperties,
				)
			}

			return schema
		}

		if unknown.IsArray() {
			schema.Properties[fullPath] = model.Property{
				Path: fullPath,
				Type: "array",
			}

			schema = buildSchema(
				deep+1,
				originalSchemaName,
				schemas,
				schema,
				newArrayPropertyPath(fullPath),
				unknown.Items,
			)

			return schema
		}

		if unknown.IsPrimitive() {
			schema.Properties[fullPath] = model.Property{
				Path: fullPath,
				Type: unknown.Type,
			}

			return schema
		}

		if unknown.IsRef() {
			schema = buildSchema(
				deep+1,
				originalSchemaName,
				schemas,
				schema,
				fullPath,
				schemas[unknown.Ref],
			)

			return schema
		}

		return schema
	case *Schema:
		return buildSchema(
			deep+1,
			originalSchemaName,
			schemas,
			schema,
			fullPath,
			*unknown,
		)
	default:
		panic(fmt.Sprintf("unknown schema type %T", unknown))
	}
}
