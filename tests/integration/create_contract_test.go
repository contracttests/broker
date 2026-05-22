package integration_test

import (
	"encoding/json"
	"net/http"

	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func (suite *Suite) TestCreateContract() {
	t := suite.T()

	const payload = `
{
	"name": "items-contract",
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

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, response.Body)

	var dslContract dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(payload), &dslContract))
	expectedContract := dslContract.ToContractModel()

	contractID := suite.DB.AssertContract("items-contract", "app")
	versionID := suite.DB.AssertContractVersion(contractID, 1, expectedContract.Checksum())
	resourceID := suite.DB.AssertResource(contractID, "provides", "rest_response", "/items", "get", "200")

	probe := model.NewProvidedRestResponse("/items", "get", "200", nil)
	probe.ContractInfo = &model.ContractInfo{Name: "items-contract", Owner: "app"}
	expectedPropertyCount := len(expectedContract.Resources[probe.PrimaryHash()].Properties)
	suite.DB.AssertPropertyCount(resourceID, expectedPropertyCount)
	suite.DB.AssertPropertyVersionChangeCounts(resourceID, versionID, expectedPropertyCount, 0)
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

	suite.DB.AssertNoContracts()
}

func (suite *Suite) TestCreateContract_ReUploadUnchanged_DoesNotCreateNewVersion() {
	t := suite.T()

	const payload = `
{
	"name": "items-contract",
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

	first, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    payload,
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
		Body:    payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, second.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, second.Body)

	contractID := suite.DB.AssertContract("items-contract", "app")
	suite.DB.AssertContractVersionCount(contractID, 1)
}

func (suite *Suite) TestCreateContract_ReUploadWithNewField_PersistsAsNewVersion() {
	t := suite.T()

	const v1Payload = `
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

	const v2Payload = `
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

	first, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    v1Payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, first.StatusCode)

	second, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    v2Payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, second.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, second.Body)

	var dslV2 dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(v2Payload), &dslV2))
	expectedV2 := dslV2.ToContractModel()

	contractID := suite.DB.AssertContract("minimal-contract", "app")
	suite.DB.AssertContractVersionCount(contractID, 2)
	v2ID := suite.DB.AssertContractVersion(contractID, 2, expectedV2.Checksum())
	suite.DB.AssertPropertyVersionCount(v2ID, 1)
	suite.DB.AssertSinglePropertyVersion(v2ID, "root.description", "added")
}

func (suite *Suite) TestCreateContract_PersistsAllResources() {
	t := suite.T()

	const payload = `
{
	"name": "items-contract",
	"owner": "app",
	"provides": {
		"rest": {
			"/items": {
				"get": {
					"responses": {"200": "Item"}
				},
				"post": {
					"request": "Item",
					"responses": {"200": "Id"}
				},
				"put": {
					"request": "Item",
					"responses": {"200": "Id"}
				},
				"delete": {
					"responses": {"204": "Id"}
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
		},
		"Id": {
			"type": "object",
			"properties": {
				"id": {"type": "string"}
			}
		}
	}
}
`

	response, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    payload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, response.Body)

	var dslContract dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(payload), &dslContract))
	expectedContract := dslContract.ToContractModel()

	contractID := suite.DB.AssertContract("items-contract", "app")
	suite.DB.AssertResourceCount(contractID, len(expectedContract.Resources))
}

func (suite *Suite) TestCreateContract_PersistsConsumedBillingInvoicesPost() {
	t := suite.T()

	const providerPayload = `
{
	"name": "billing",
	"owner": "billing-team",
	"provides": {
		"rest": {
			"/invoices": {
				"post": {
					"request": "Item",
					"responses": {"201": "Id"}
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
		},
		"Id": {
			"type": "object",
			"properties": {
				"id": {"type": "string"}
			}
		}
	}
}
`

	const consumerPayload = `
{
	"name": "items-contract",
	"owner": "app",
	"consumes": {
		"billing": {
			"rest": {
				"/invoices": {
					"post": {
						"request": "Item",
						"responses": {"201": "Id"}
					}
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
		},
		"Id": {
			"type": "object",
			"properties": {
				"id": {"type": "string"}
			}
		}
	}
}
`

	provider, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    providerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to seed billing: %v", err)
	}
	assert.Equal(t, http.StatusOK, provider.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, provider.Body)

	consumer, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    consumerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Fatalf("Failed to upload consumer: %v", err)
	}
	assert.Equal(t, http.StatusOK, consumer.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, consumer.Body)

	contractID := suite.DB.AssertContract("items-contract", "app")
	suite.DB.AssertResource(contractID, "consumes", "rest_response", "/invoices", "post", "201")
}
