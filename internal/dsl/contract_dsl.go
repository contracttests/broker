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
	resources := c.buildResources([]model.Resource{}, ResourcePath(""), *c)
	return model.Contract{
		Resources: resources,
	}
}

func (c *Contract) buildResources(resources []model.Resource, resourcePath ResourcePath, unknown any) []model.Resource {
	switch unknown := unknown.(type) {
	case Contract:
		resourcePath = resourcePath.Append(unknown.Api.Name)

		for serviceName, consumes := range unknown.ConsumesServices {
			newResourcePath := resourcePath.Append("consumes", serviceName)
			resources = c.buildResources(resources, newResourcePath, consumes)
		}

		newResourcePath := resourcePath.Append("provides")
		resources = c.buildResources(resources, newResourcePath, unknown.Provides)

		return resources

	case Consumes:
		resources = c.buildResources(resources, resourcePath, unknown.Rest)
		resources = c.buildResources(resources, resourcePath, unknown.Message)

		return resources

	case Provides:
		resources = c.buildResources(resources, resourcePath, unknown.Rest)
		resources = c.buildResources(resources, resourcePath, unknown.Message)

		return resources

	case Message:
		for messageName, schemaName := range unknown {
			newResourcePath := resourcePath.Append("message", messageName)
			schema := buildSchema(
				0,
				schemaName,
				c.Schemas,
				model.Schema{
					Properties: make(map[string]model.Property),
				},
				PropertyPath("root"),
				c.Schemas[schemaName],
			)
			resources = append(resources, model.NewResource(newResourcePath.String(), schema))
		}

		return resources

	case Rest:
		for endpoint, methods := range unknown {
			if methods.Get.IsNonZero() {
				newResourcePath := resourcePath.Append("rest", endpoint)
				resources = c.buildResources(resources, newResourcePath, methods.Get)
			}

			if methods.Post.IsNonZero() {
				newResourcePath := resourcePath.Append("rest", endpoint)
				resources = c.buildResources(resources, newResourcePath, methods.Post)
			}

			if methods.Put.IsNonZero() {
				newResourcePath := resourcePath.Append("rest", endpoint)
				resources = c.buildResources(resources, newResourcePath, methods.Put)
			}

			if methods.Delete.IsNonZero() {
				newResourcePath := resourcePath.Append("rest", endpoint)
				resources = c.buildResources(resources, newResourcePath, methods.Delete)
			}
		}

	case GetMethod:
		newResourcePath := resourcePath.Append("get", "responses")
		resources = c.buildResources(resources, newResourcePath, unknown.Responses)

		return resources

	case PostMethod:
		if unknown.HasRequestBody() {
			newResourcePath := resourcePath.Append("post", "request")

			schema := buildSchema(
				0,
				unknown.RequestBody,
				c.Schemas,
				model.Schema{
					Properties: make(map[string]model.Property),
				},
				PropertyPath("root"),
				c.Schemas[unknown.RequestBody],
			)

			resources = append(resources, model.NewResource(newResourcePath.String(), schema))
		}

		newResourcePath := resourcePath.Append("post", "responses")
		resources = c.buildResources(resources, newResourcePath, unknown.Responses)

		return resources

	case PutMethod:
		if unknown.HasRequestBody() {
			newResourcePath := resourcePath.Append("put", "request")

			schema := buildSchema(
				0,
				unknown.RequestBody,
				c.Schemas,
				model.Schema{
					Properties: make(map[string]model.Property),
				},
				PropertyPath("root"),
				c.Schemas[unknown.RequestBody],
			)

			resources = append(resources, model.NewResource(newResourcePath.String(), schema))
		}

		newResourcePath := resourcePath.Append("put", "responses")
		resources = c.buildResources(resources, newResourcePath, unknown.Responses)

		return resources

	case DeleteMethod:
		newResourcePath := resourcePath.Append("delete", "responses")
		resources = c.buildResources(resources, newResourcePath, unknown.Responses)

		return resources

	case Responses:
		for statusCode, schemaName := range unknown {
			newResourcePath := resourcePath.Append(strconv.Itoa(statusCode))
			schema := buildSchema(
				0,
				schemaName,
				c.Schemas,
				model.Schema{
					Properties: make(map[string]model.Property),
				},
				PropertyPath("root"),
				c.Schemas[schemaName],
			)

			resources = append(resources, model.NewResource(newResourcePath.String(), schema))
		}

		return resources
	}

	return resources
}

func (c *Contract) SchemasResolver(contractDsl Contract) map[string]model.Schema {
	schemas := make(map[string]model.Schema)

	for schemaName, schema := range contractDsl.Schemas {
		hash := model.UuidFromStrings(schemaName)

		schema := buildSchema(
			0,
			schemaName,
			contractDsl.Schemas,
			model.Schema{
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
