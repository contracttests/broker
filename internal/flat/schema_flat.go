package flat

import (
	"fmt"
	"strings"

	"github.com/contracttests/broker/internal/dsl"
)

type FlatSchemas map[string]FlatSchema

type FlatSchema []FlatProperty

type FlatProperty struct {
	FullPath string
	Type     string
}

func (f *FlatProperty) IsPrimitive() bool {
	return f.Type == "string" ||
		f.Type == "number" ||
		f.Type == "boolean" ||
		f.Type == "integer" ||
		f.Type == "object" ||
		f.Type == "array"
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

func Schemas(contractDsl dsl.Contract) FlatSchemas {
	flatSchemas := FlatSchemas{}

	for schemaName, schema := range contractDsl.Schemas {
		flatSchema := buildFlatProperties(
			0,
			schemaName,
			contractDsl.Schemas,
			FlatSchema{},
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
	schemas dsl.Schemas,
	flatSchema FlatSchema,
	fullPath string,
	unknown any,
) FlatSchema {
	if deep >= 10 {
		panic(fmt.Sprintf("Circular reference detected in the schema %s", originalSchemaName))
	}

	switch unknown := unknown.(type) {
	case dsl.Schema:
		if unknown.IsObject() {
			flatSchema = append(flatSchema, FlatProperty{
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
			flatSchema = append(flatSchema, FlatProperty{
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
			flatSchema = append(flatSchema, FlatProperty{
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
	case *dsl.Schema:
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
