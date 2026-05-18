package integration_test

import (
	"encoding/json"
	"net/http"

	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/model"
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

func (suite *Suite) TestCreateContract() {
	t := suite.T()

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    itemsContractPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, response.Body)

	var dslContract dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(itemsContractPayload), &dslContract))
	expectedContract := dslContract.ToContractModel()

	contractID := suite.AssertContract("items-contract", "app")
	versionID := suite.AssertContractVersion(contractID, 1, expectedContract.Checksum())
	resourceID := suite.AssertResource(contractID, "provides", "rest_response", "/items", "get", "200")

	expectedPropertyCount := len(expectedContract.Resources[model.NewProvidedRestResponse("/items", "get", "200", nil).Key()].Properties)
	suite.AssertPropertyCount(resourceID, expectedPropertyCount)
	suite.AssertPropertyVersionChangeCounts(resourceID, versionID, expectedPropertyCount, 0)
}

func (suite *Suite) TestCreateContract_MissingName_Returns400() {
	t := suite.T()

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    `{"owner": "app"}`,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.JSONEq(t, `{"success":false,"message":"contract invalid input"}`, response.Body)

	suite.AssertNoContracts()
}

func (suite *Suite) TestCreateContract_ReUploadUnchanged_DoesNotCreateNewVersion() {
	t := suite.T()

	first, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    itemsContractPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, first.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, first.Body)

	second, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    itemsContractPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, second.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, second.Body)

	contractID := suite.AssertContract("items-contract", "app")
	suite.AssertContractVersionCount(contractID, 1)
}

const minimalContractV1Payload = `
{
	"name": "minimal-contract",
	"owner": "app",
	"provides": {
		"rest": {
			"/items": {
				"get": {
					"responses": {"200": "Item"}
				}
			}
		}
	},
	"schemas": {
		"Item": {
			"type": "object",
			"properties": {
				"name": {"type": "string"}
			}
		}
	}
}
`

const minimalContractV2Payload = `
{
	"name": "minimal-contract",
	"owner": "app",
	"provides": {
		"rest": {
			"/items": {
				"get": {
					"responses": {"200": "Item"}
				}
			}
		}
	},
	"schemas": {
		"Item": {
			"type": "object",
			"properties": {
				"name": {"type": "string"},
				"description": {"type": "string"}
			}
		}
	}
}
`

func (suite *Suite) TestCreateContract_ReUploadWithNewField_PersistsAsNewVersion() {
	t := suite.T()

	first, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    minimalContractV1Payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, first.StatusCode)

	second, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    minimalContractV2Payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, second.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, second.Body)

	var dslV2 dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(minimalContractV2Payload), &dslV2))
	expectedV2 := dslV2.ToContractModel()

	contractID := suite.AssertContract("minimal-contract", "app")
	suite.AssertContractVersionCount(contractID, 2)
	v2ID := suite.AssertContractVersion(contractID, 2, expectedV2.Checksum())
	suite.AssertPropertyVersionCount(v2ID, 1)
	suite.AssertSinglePropertyVersion(v2ID, "root.description", "added")
}

func (suite *Suite) TestCreateContract_PersistsAllResources() {
	t := suite.T()

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    itemsContractPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var dslContract dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(itemsContractPayload), &dslContract))
	expectedContract := dslContract.ToContractModel()

	contractID := suite.AssertContract("items-contract", "app")
	suite.AssertResourceCount(contractID, len(expectedContract.Resources))
}

func (suite *Suite) TestCreateContract_PersistsConsumedBillingInvoicesPost() {
	t := suite.T()

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    itemsContractPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contractID := suite.AssertContract("items-contract", "app")
	suite.AssertResource(contractID, "consumes", "rest_response", "/invoices", "post", "201")
}
