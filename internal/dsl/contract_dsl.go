package dsl

import (
	"fmt"
	"strconv"

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

func (c *Contract) ToContractModel() model.Contract {
	resources := c.buildResources([]model.Resource{}, NewResourcePath(""), *c)
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
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[schemaName],
			)
			resources = append(resources, NewResource(newResourcePath, schema))
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
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[unknown.RequestBody],
			)

			resources = append(resources, NewResource(newResourcePath, schema))
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
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[unknown.RequestBody],
			)

			resources = append(resources, NewResource(newResourcePath, schema))
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
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[schemaName],
			)

			resources = append(resources, NewResource(newResourcePath, schema))
		}

		return resources
	}

	return resources
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
			schema.Properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), "object")

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
			schema.Properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), "array")

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
			schema.Properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), unknown.Type)

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
