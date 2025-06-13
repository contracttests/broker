package dsl_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestResouce(t *testing.T) {
	content := []byte(`
		{
			"get": {
				"responses": {
					"200": "Products"
				}
			},
			"post": {
				"requestBody": "CreateProduct",
				"responses": {
					"201": "Product",
					"400": "BadRequest"
				}
			},
			"put": {
				"requestBody": "UpdateProduct",
				"responses": {
					"200": "Product",
					"400": "BadRequest",
					"404": "NotFound"
				}
			},
			"delete": {
				"responses": {
					"204": "NoContent",
					"404": "NotFound"
				}
			}
		}
	`)

	actual := dsl.Resource{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.Resource{
		Get: dsl.GetMethod{
			Responses: dsl.Responses{
				200: "Products",
			},
		},
		Post: dsl.PostMethod{
			RequestBody: "CreateProduct",
			Responses: dsl.Responses{
				201: "Product",
				400: "BadRequest",
			},
		},
		Put: dsl.PutMethod{
			RequestBody: "UpdateProduct",
			Responses: dsl.Responses{
				200: "Product",
				400: "BadRequest",
				404: "NotFound",
			},
		},
		Delete: dsl.DeleteMethod{
			Responses: dsl.Responses{
				204: "NoContent",
				404: "NotFound",
			},
		},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
}

func TestZeroResource(t *testing.T) {
	content := []byte(`
		{
			"get": {},
			"post": {},
			"put": {},
			"delete": {}
		}
	`)

	actual := dsl.Resource{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.Resource{
		Get:    dsl.GetMethod{},
		Post:   dsl.PostMethod{},
		Put:    dsl.PutMethod{},
		Delete: dsl.DeleteMethod{},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
}
