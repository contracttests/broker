package upload_contract_test

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/assert"
)

const itemsContractPayload = `
{
	"name": "items-contract",
	"owner": "app",
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

func (suite *Suite) postItemsContract() *Response {
	response, err := suite.Request(Request{
		Method: "POST",
		Path:   "/contracts",
		Body:   itemsContractPayload,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err != nil {
		suite.T().Fatalf("Failed to test request: %v", err)
	}

	return response
}

func (suite *Suite) TestCreateContract() {
	t := suite.T()

	response := suite.postItemsContract()

	assert.Equal(t, http.StatusOK, response.StatusCode)

	expected := `
	{
		"success":true,
		"message":"contract upload successful"
	}`
	assert.JSONEq(t, expected, response.Body)
}

func (suite *Suite) TestCreateContract_PersistsFullTree() {
	t := suite.T()
	ctx := context.Background()

	response := suite.postItemsContract()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var contractID int64
	var owner string
	err := suite.Components.Pool.QueryRow(ctx,
		`SELECT id, owner FROM contracts WHERE name = $1`, "items-contract",
	).Scan(&contractID, &owner)
	assert.NoError(t, err)
	assert.Equal(t, "app", owner)

	var versionCount int
	var checksum string
	var rawPayloadLen int
	err = suite.Components.Pool.QueryRow(ctx,
		`SELECT COUNT(*), MAX(checksum), MAX(LENGTH(raw_payload::text))
		 FROM contract_versions WHERE contract_id = $1 AND version = 1`, contractID,
	).Scan(&versionCount, &checksum, &rawPayloadLen)
	assert.NoError(t, err)
	assert.Equal(t, 1, versionCount)
	assert.NotEmpty(t, checksum)
	assert.Positive(t, rawPayloadLen)

	var resourceCount int
	err = suite.Components.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM resources WHERE contract_id = $1`, contractID,
	).Scan(&resourceCount)
	assert.NoError(t, err)
	assert.Positive(t, resourceCount)

	var propertyCount int
	err = suite.Components.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM properties p
		 JOIN resources r ON r.id = p.resource_id
		 WHERE r.contract_id = $1`, contractID,
	).Scan(&propertyCount)
	assert.NoError(t, err)
	assert.Positive(t, propertyCount)

	var addedCount, nonAddedCount int
	err = suite.Components.Pool.QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE pv.change = 'added'),
			COUNT(*) FILTER (WHERE pv.change != 'added')
		 FROM property_versions pv
		 JOIN properties p ON p.id = pv.property_id
		 JOIN resources r ON r.id = p.resource_id
		 WHERE r.contract_id = $1`, contractID,
	).Scan(&addedCount, &nonAddedCount)
	assert.NoError(t, err)
	assert.Equal(t, propertyCount, addedCount)
	assert.Equal(t, 0, nonAddedCount)
}

func (suite *Suite) TestCreateContract_MissingName_Returns400() {
	t := suite.T()

	payload := `{"owner": "app"}`

	response, err := suite.Request(Request{
		Method: "POST",
		Path:   "/contracts",
		Body:   payload,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	expected := `
	{
		"success":false,
		"message":"contract invalid input"
	}`
	assert.JSONEq(t, expected, response.Body)
}

func (suite *Suite) TestCreateContract_DuplicateName_ReturnsAlreadyUploaded() {
	t := suite.T()
	ctx := context.Background()

	first := suite.postItemsContract()
	assert.Equal(t, http.StatusOK, first.StatusCode)

	second := suite.postItemsContract()
	assert.Equal(t, http.StatusOK, second.StatusCode)

	expected := `
	{
		"success":true,
		"message":"contract already uploaded"
	}`
	assert.JSONEq(t, expected, second.Body)

	var count int
	err := suite.Components.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM contracts WHERE name = $1`, "items-contract",
	).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
