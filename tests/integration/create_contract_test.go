package upload_contract_test

import (
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (suite *Suite) TestCreateContract() {
	t := suite.T()

	payload := `
	{
		"api": {
			"name": "app"
		},
		"provides": {
			"rest": {
				"/items": {
					"get": {
						"responses": {
							"200": "Item"
						}
					},
					"post": {
						"request": "Item",
						"responses": {
							"200": "Id"
						}
					},
					"put": {
						"request": "Item",
						"responses": {
							"200": "Id"
						}
					},
					"delete": {
						"responses": {
							"204": "Id"
						}
					}
				}
			}
		},
		"consumes": {
			"billing": {
				"rest": {
					"/invoices": {
						"get": {
							"responses": {
								"200": "Item"
							}
						},
						"post": {
							"request": "Item",
							"responses": {
								"201": "Id"
							}
						},
						"put": {
							"request": "Item",
							"responses": {
								"200": "Id"
							}
						},
						"delete": {
							"responses": {
								"200": "Id"
							}
						}
					}
				}
			}
		},
		"schemas": {
			"Item": {
			"type": "object",
			"properties": {
				"name": {
					"type": "string"
				},
				"tags": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"owner": {
					"$ref": "User"
				}
			}
			},
			"User": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string"
					}
				}
			},
			"Id": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string"
					}
				}
			}
		}
	}
	`

	response, err := suite.Request(Request{
		Method: "POST",
		Path: "/contracts",
		Body: payload,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	expected := `
	{
		"success":true,
		"message":"contract upload successful"
	}`
	assert.JSONEq(t, expected, response.Body)
}
