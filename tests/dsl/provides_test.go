package dsl_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestProvides(t *testing.T) {
	content := []byte(`
		{
			"rest": {
				"/products": {
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
			}
		}
	`)

	actual := dsl.Provides{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.Provides{
		Rest: map[string]dsl.Resource{
			"/products": {
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
			},
		},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
}
