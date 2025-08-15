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

	if string(*f) == "" {
		return ResourcePath(strings.Join(parts, separator))
	}

	return ResourcePath(strings.Join([]string{string(*f), strings.Join(parts, separator)}, separator))
}

func (f *ResourcePath) String() string {
	return string(*f)
}

type PropertyPath string

func (f *PropertyPath) Append(parts ...string) PropertyPath {
	separator := "."

	if string(*f) == "" {
		return PropertyPath(strings.Join(parts, separator))
	}

	return PropertyPath(strings.Join([]string{string(*f), strings.Join(parts, separator)}, separator))
}

func (f *PropertyPath) AppendArray() PropertyPath {
	return PropertyPath(f.String() + "[]")
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

func Resources(contractDsl Contract) []model.Resource {
	return buildResources([]model.Resource{}, ResourcePath(""), contractDsl)
}

func buildResources(flatResources []model.Resource, resourcePath ResourcePath, unknown any) []model.Resource {
	switch unknown := unknown.(type) {
	case Contract:
		resourcePath = resourcePath.Append(unknown.Api.Name)

		for serviceName, consumes := range unknown.ConsumesServices {
			newFullPath := resourcePath.Append("consumes", serviceName)
			flatResources = buildResources(flatResources, newFullPath, consumes)
		}

		newFullPath := resourcePath.Append("provides")
		flatResources = buildResources(flatResources, newFullPath, unknown.Provides)

		return flatResources

	case Consumes:
		flatResources = buildResources(flatResources, resourcePath, unknown.Rest)
		flatResources = buildResources(flatResources, resourcePath, unknown.Message)

		return flatResources

	case Provides:
		flatResources = buildResources(flatResources, resourcePath, unknown.Rest)
		flatResources = buildResources(flatResources, resourcePath, unknown.Message)

		return flatResources

	case Message:
		for messageName, schemaName := range unknown {
			newFullPath := resourcePath.Append("message", messageName)
			flatResources = append(flatResources, model.NewResource(newFullPath.String(), model.UuidFromStrings(schemaName)))
		}

		return flatResources

	case Rest:
		for endpoint, methods := range unknown {
			if methods.Get.IsNonZero() {
				newFullPath := resourcePath.Append("rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Get)
			}

			if methods.Post.IsNonZero() {
				newFullPath := resourcePath.Append("rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Post)
			}

			if methods.Put.IsNonZero() {
				newFullPath := resourcePath.Append("rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Put)
			}

			if methods.Delete.IsNonZero() {
				newFullPath := resourcePath.Append("rest", endpoint)
				flatResources = buildResources(flatResources, newFullPath, methods.Delete)
			}
		}

	case GetMethod:
		newFullPath := resourcePath.Append("get", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PostMethod:
		if unknown.HasRequestBody() {
			newFullPath := resourcePath.Append("post", "request")
			flatResources = append(flatResources, model.NewResource(newFullPath.String(), model.UuidFromStrings(unknown.RequestBody)))
		}

		newFullPath := resourcePath.Append("post", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PutMethod:
		if unknown.HasRequestBody() {
			newFullPath := resourcePath.Append("put", "request")
			flatResources = append(flatResources, model.NewResource(newFullPath.String(), model.UuidFromStrings(unknown.RequestBody)))
		}

		newFullPath := resourcePath.Append("put", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case DeleteMethod:
		newFullPath := resourcePath.Append("delete", "responses")
		flatResources = buildResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case Responses:
		for statusCode, schemaName := range unknown {
			newFullPath := resourcePath.Append(strconv.Itoa(statusCode))
			flatResources = append(flatResources, model.NewResource(newFullPath.String(), model.UuidFromStrings(schemaName)))
		}

		return flatResources
	}

	return flatResources
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
			PropertyPath("root"),
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
	propertyPath PropertyPath,
	unknown any,
) model.Schema {
	if deep >= 10 {
		panic(fmt.Sprintf("Circular reference detected in the schema %s", originalSchemaName))
	}

	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			schema.Properties[propertyPath.String()] = model.Property{
				Path: propertyPath.String(),
				Type: "object",
			}

			for name, schemaProperties := range unknown.Properties {
				schema = buildSchema(
					deep+1,
					originalSchemaName,
					schemas,
					schema,
					propertyPath.Append(name),
					schemaProperties,
				)
			}

			return schema
		}

		if unknown.IsArray() {
			schema.Properties[propertyPath.String()] = model.Property{
				Path: propertyPath.String(),
				Type: "array",
			}

			schema = buildSchema(
				deep+1,
				originalSchemaName,
				schemas,
				schema,
				propertyPath.AppendArray(),
				unknown.Items,
			)

			return schema
		}

		if unknown.IsPrimitive() {
			schema.Properties[propertyPath.String()] = model.Property{
				Path: propertyPath.String(),
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
				propertyPath,
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
			propertyPath,
			*unknown,
		)
	default:
		panic(fmt.Sprintf("unknown schema type %T", unknown))
	}
}
