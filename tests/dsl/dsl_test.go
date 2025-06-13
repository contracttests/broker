package dsl_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestPrimitiveType(t *testing.T) {
	content := []byte(`
		{ 
			"Product": {
				"type":"object",
				"properties":{
					"uuid":{
						"type":"string"
					},
					"name":{
						"type":"string"
					},
					"birthdate":{
						"type":"string"
					}
				}
			}
		 }
	`)

	schema := dsl.Schemas{}

	if err := json.Unmarshal(content, &schema); err != nil {
		log.Fatal("could not parse json, make sure it is valid contract file")
	}

	actual := dsl.NewComparableSchemas(schema)

	root := dsl.ComparableSchemaProperty{
		Path: "root",
		Type: "object",
	}

	uuid := dsl.ComparableSchemaProperty{
		Path: "root.uuid",
		Type: "string",
	}

	name := dsl.ComparableSchemaProperty{
		Path: "root.name",
		Type: "string",
	}
	birthdate := dsl.ComparableSchemaProperty{
		Path: "root.birthdate",
		Type: "string",
	}

	properties := dsl.ComparableSchemaProperties{
		root,
		uuid,
		name,
		birthdate,
	}

	expected := dsl.ComparableSchemas{
		dsl.ComparableSchema{
			Name:       "Product",
			Properties: properties,
		},
	}

	assert.ElementsMatch(t, expected, actual)
}
