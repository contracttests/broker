package integration_test

import (
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (suite *Suite) TestUploadConsumer_MissingTwoProperties_Returns422WithBothBreaks() {
	t := suite.T()

	const providerPayload = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets": {
				"get": {"responses": {"200": "Pets"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"},
				"name": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const consumerPayload = `
{
	"name": "broken-app",
	"owner": "app-team",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets": {
					"get": {"responses": {"200": "Pets"}}
				}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"},
				"name": {"type": "string"},
				"deletedAt": {"type": "string"},
				"soldAt": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	providerResp, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    providerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusOK, providerResp.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, providerResp.Body)

	consumerResp, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    consumerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, consumerResp.StatusCode)
	assert.JSONEq(t, `{
		"success": false,
		"message": "contract incompatible with stored counterparts",
		"breakingChanges": [
			{
				"contractName": "broken-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].deletedAt",
				"reason": "missing_in_provider"
			},
			{
				"contractName": "broken-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].soldAt",
				"reason": "missing_in_provider"
			}
		]
	}`, consumerResp.Body)

	suite.DB.AssertContractNotPersisted("broken-app")
}

func (suite *Suite) TestUploadConsumer_TwoUnknownProviders_Returns422WithBothBreaks() {
	t := suite.T()

	const consumerPayload = `
{
	"name": "lonely-app",
	"owner": "app-team",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets": {
					"get": {"responses": {"200": "Pets"}}
				}
			}
		},
		"users-service": {
			"rest": {
				"/users": {
					"get": {"responses": {"200": "Users"}}
				}
			}
		}
	},
	"schemas": {
		"Pet":   {"type": "object", "properties": {"uuid": {"type": "string"}}},
		"Pets":  {"type": "array", "items": {"$ref": "Pet"}},
		"User":  {"type": "object", "properties": {"uuid": {"type": "string"}}},
		"Users": {"type": "array", "items": {"$ref": "User"}}
	}
}`

	resp, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    consumerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.JSONEq(t, `{
		"success": false,
		"message": "contract incompatible with stored counterparts",
		"breakingChanges": [
			{
				"contractName": "lonely-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"reason": "provider_resource_not_found"
			},
			{
				"contractName": "lonely-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"users-service","endpoint":"/users","method":"get","statusCode":"200"},
				"reason": "provider_resource_not_found"
			}
		]
	}`, resp.Body)

	suite.DB.AssertContractNotPersisted("lonely-app")
}

func (suite *Suite) TestUploadConsumer_TwoTypeMismatches_Returns422CarriesExpectedAndActualTypes() {
	t := suite.T()

	const providerPayload = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets/{uuid}": {
				"get": {"responses": {"200": "Pet"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"age":   {"type": "integer"},
				"score": {"type": "integer"}
			}
		}
	}
}`

	const consumerPayload = `
{
	"name": "mismatched-app",
	"owner": "app-team",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets/{uuid}": {
					"get": {"responses": {"200": "Pet"}}
				}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"age":   {"type": "string"},
				"score": {"type": "string"}
			}
		}
	}
}`

	providerResp, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    providerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusOK, providerResp.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, providerResp.Body)

	consumerResp, err := suite.Request(Request{
		Method:  "POST",
		Path:    "/contracts",
		Body:    consumerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, consumerResp.StatusCode)
	assert.JSONEq(t, `{
		"success": false,
		"message": "contract incompatible with stored counterparts",
		"breakingChanges": [
			{
				"contractName": "mismatched-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets/{uuid}","method":"get","statusCode":"200"},
				"property": "root.age",
				"reason": "type_mismatch",
				"expectedType": "string",
				"actualType": "integer"
			},
			{
				"contractName": "mismatched-app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets/{uuid}","method":"get","statusCode":"200"},
				"property": "root.score",
				"reason": "type_mismatch",
				"expectedType": "string",
				"actualType": "integer"
			}
		]
	}`, consumerResp.Body)
}

func (suite *Suite) TestUploadProvider_DropsTwoPropertiesNeededByConsumer_Returns422AndDoesNotUpdate() {
	t := suite.T()

	const providerV1 = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets": {
				"get": {"responses": {"200": "Pets"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid":  {"type": "string"},
				"name":  {"type": "string"},
				"email": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const consumerPayload = `
{
	"name": "app",
	"owner": "app-team",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets": {
					"get": {"responses": {"200": "Pets"}}
				}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid":  {"type": "string"},
				"name":  {"type": "string"},
				"email": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const providerV2 = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets": {
				"get": {"responses": {"200": "Pets"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	v1, err := suite.Request(Request{
		Method: "POST", Path: "/contracts", Body: providerV1,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusOK, v1.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, v1.Body)

	consumer, err := suite.Request(Request{
		Method: "POST", Path: "/contracts", Body: consumerPayload,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusOK, consumer.StatusCode)
	assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, consumer.Body)

	v2, err := suite.Request(Request{
		Method: "POST", Path: "/contracts", Body: providerV2,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, v2.StatusCode)
	assert.JSONEq(t, `{
		"success": false,
		"message": "contract incompatible with stored counterparts",
		"breakingChanges": [
			{
				"contractName": "app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].email",
				"reason": "missing_in_provider"
			},
			{
				"contractName": "app",
				"contractOwner": "app-team",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].name",
				"reason": "missing_in_provider"
			}
		]
	}`, v2.Body)

	providerContractID := suite.DB.AssertContract("pets-service", "pets-team")
	suite.DB.AssertContractVersionCount(providerContractID, 1)
}

func (suite *Suite) TestUploadProvider_DropsPropertyNeededByTwoConsumers_Returns422WithBreakPerConsumer() {
	t := suite.T()

	const providerV1 = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets": {
				"get": {"responses": {"200": "Pets"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"},
				"name": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const consumerA = `
{
	"name": "app-a",
	"owner": "team-a",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets": {
					"get": {"responses": {"200": "Pets"}}
				}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"},
				"name": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const consumerB = `
{
	"name": "app-b",
	"owner": "team-b",
	"consumes": {
		"pets-service": {
			"rest": {
				"/pets": {
					"get": {"responses": {"200": "Pets"}}
				}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"},
				"name": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	const providerV2 = `
{
	"name": "pets-service",
	"owner": "pets-team",
	"provides": {
		"rest": {
			"/pets": {
				"get": {"responses": {"200": "Pets"}}
			}
		}
	},
	"schemas": {
		"Pet": {
			"type": "object",
			"properties": {
				"uuid": {"type": "string"}
			}
		},
		"Pets": {"type": "array", "items": {"$ref": "Pet"}}
	}
}`

	for _, payload := range []string{providerV1, consumerA, consumerB} {
		resp, err := suite.Request(Request{
			Method: "POST", Path: "/contracts", Body: payload,
			Headers: map[string]string{"Content-Type": "application/json"},
		})
		suite.Require().NoError(err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, `{"success":true,"message":"contract upload successful"}`, resp.Body)
	}

	resp, err := suite.Request(Request{
		Method: "POST", Path: "/contracts", Body: providerV2,
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	suite.Require().NoError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.JSONEq(t, `{
		"success": false,
		"message": "contract incompatible with stored counterparts",
		"breakingChanges": [
			{
				"contractName": "app-a",
				"contractOwner": "team-a",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].name",
				"reason": "missing_in_provider"
			},
			{
				"contractName": "app-b",
				"contractOwner": "team-b",
				"resource": {"direction":"consumes","kind":"rest_response","provider":"pets-service","endpoint":"/pets","method":"get","statusCode":"200"},
				"property": "root[].name",
				"reason": "missing_in_provider"
			}
		]
	}`, resp.Body)
}
