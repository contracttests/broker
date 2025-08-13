package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/contracttests/broker/internal/flat"
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

func (c *Contract) ToFlatContract() flat.FlatContract {
	flatResources := Resources(*c)
	flatSchemas := Schemas(*c)

	return flat.FlatContract{
		Resources: flatResources,
		Schemas:   flatSchemas,
	}
}

func newResourceFullPath(parts ...string) string {
	resourceParthSeparator := ";"

	return strings.Join(parts, resourceParthSeparator)
}

func Resources(contractDsl Contract) []flat.FlatResource {
	return buildFlatResources([]flat.FlatResource{}, "", contractDsl)
}

func buildFlatResources(flatResources []flat.FlatResource, fullPath string, unknown any) []flat.FlatResource {
	switch unknown := unknown.(type) {
	case Contract:
		fullPath = unknown.Api.Name

		for serviceName, consumes := range unknown.ConsumesServices {
			newFullPath := newResourceFullPath(fullPath, "consumes", serviceName)
			flatResources = buildFlatResources(flatResources, newFullPath, consumes)
		}

		newFullPath := newResourceFullPath(fullPath, "provides")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Provides)

		return flatResources

	case Consumes:
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case Provides:
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case Message:
		for messageName, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, "message", messageName)
			flatResources = append(flatResources, flat.FlatResource{
				FullPath:   newFullPath,
				SchemaName: schemaName,
			})
		}

		return flatResources

	case Rest:
		for endpoint, methods := range unknown {
			if methods.Get.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Get)
			}

			if methods.Post.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Post)
			}

			if methods.Put.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Put)
			}

			if methods.Delete.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Delete)
			}
		}

	case GetMethod:
		newFullPath := newResourceFullPath(fullPath, "get", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PostMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "post", "requestBody")
			flatResources = append(flatResources, flat.FlatResource{
				FullPath:   newFullPath,
				SchemaName: unknown.RequestBody,
			})
		}

		newFullPath := newResourceFullPath(fullPath, "post", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case PutMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "put", "requestBody")
			flatResources = append(flatResources, flat.FlatResource{
				FullPath:   newFullPath,
				SchemaName: unknown.RequestBody,
			})
		}

		newFullPath := newResourceFullPath(fullPath, "put", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case DeleteMethod:
		newFullPath := newResourceFullPath(fullPath, "delete", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case Responses:
		for statusCode, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, strconv.Itoa(statusCode))
			flatResources = append(flatResources, flat.FlatResource{
				FullPath:   newFullPath,
				SchemaName: schemaName,
			})
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

func Schemas(contractDsl Contract) flat.FlatSchemas {
	flatSchemas := flat.FlatSchemas{}

	for schemaName, schema := range contractDsl.Schemas {
		flatSchema := buildFlatProperties(
			0,
			schemaName,
			contractDsl.Schemas,
			flat.FlatSchema{},
			"root",
			schema,
		)
		flatSchemas[schemaName] = flatSchema
	}

	return flatSchemas
}

func buildFlatProperties(
	deep int,
	originalSchemaName string,
	schemas map[string]Schema,
	flatSchema flat.FlatSchema,
	fullPath string,
	unknown any,
) flat.FlatSchema {
	if deep >= 10 {
		panic(fmt.Sprintf("Circular reference detected in the schema %s", originalSchemaName))
	}

	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			flatSchema = append(flatSchema, flat.FlatProperty{
				FullPath: fullPath,
				Type:     "object",
			})

			for name, schema := range unknown.Properties {
				flatSchema = buildFlatProperties(
					deep+1,
					originalSchemaName,
					schemas,
					flatSchema,
					newFullPath(fullPath, name),
					schema,
				)
			}

			return flatSchema
		}

		if unknown.IsArray() {
			flatSchema = append(flatSchema, flat.FlatProperty{
				FullPath: fullPath,
				Type:     "array",
			})

			flatSchema = buildFlatProperties(
				deep+1,
				originalSchemaName,
				schemas,
				flatSchema,
				newArrayPropertyPath(fullPath),
				unknown.Items,
			)

			return flatSchema
		}

		if unknown.IsPrimitive() {
			flatSchema = append(flatSchema, flat.FlatProperty{
				FullPath: newFullPath(fullPath),
				Type:     unknown.Type,
			})

			return flatSchema
		}

		if unknown.IsRef() {
			flatSchema = buildFlatProperties(
				deep+1,
				originalSchemaName,
				schemas,
				flatSchema,
				fullPath,
				schemas[unknown.Ref],
			)

			return flatSchema
		}

		return flatSchema
	case *Schema:
		return buildFlatProperties(
			deep+1,
			originalSchemaName,
			schemas,
			flatSchema,
			fullPath,
			*unknown,
		)
	default:
		panic(fmt.Sprintf("unknown schema type %T", unknown))
	}
}
