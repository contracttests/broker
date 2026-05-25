package dsl_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const happyContractJSON = `{
  "consumes": {
    "pets-service": {
      "rest": {
        "/pets": {
          "get": {
            "responses": {
              "200": "Pet"
            }
          }
        }
      }
    }
  },
  "schemas": {
    "Pet": {
      "type": "object",
      "properties": {
        "id":   { "type": "string" },
        "name": { "type": "string" }
      }
    }
  }
}`

const cyclicContractJSON = `{
  "provides": {
    "rest": {
      "/pets": {
        "get": {
          "responses": {
            "200": "Pet"
          }
        }
      }
    }
  },
  "schemas": {
    "Pet": {
      "type": "object",
      "properties": {
        "self": { "$ref": "Pet" }
      }
    }
  }
}`

func TestHydrateContract_Happy_MaterializesResources(t *testing.T) {
	var dslContract dsl.Contract
	require.NoError(t, json.Unmarshal([]byte(happyContractJSON), &dslContract))

	contract := model.NewContract(model.NewParticipant("petstore-app"), happyContractJSON)
	dslContract.HydrateContract(contract)

	require.Len(t, contract.Resources, 1)

	var resource model.Resource
	for _, r := range contract.Resources {
		resource = r
	}

	assert.Equal(t, model.Consumes, resource.Direction)
	assert.Equal(t, model.RestResponse, resource.Kind)
	assert.Equal(t, "pets-service", resource.Provider)
	assert.Equal(t, "/pets", resource.Endpoint)
	assert.Equal(t, "get", resource.Method)
	assert.Equal(t, "200", resource.StatusCode)

	assert.Contains(t, resource.Properties, "root")
	assert.Contains(t, resource.Properties, "root.id")
	assert.Contains(t, resource.Properties, "root.name")
	assert.Equal(t, "string", resource.Properties["root.id"].Type)
}

func TestHydrateContract_Unhappy_PanicsOnSchemaTooDeep(t *testing.T) {
	var dslContract dsl.Contract
	require.NoError(t, json.Unmarshal([]byte(cyclicContractJSON), &dslContract))

	contract := model.NewContract(model.NewParticipant("petstore-app"), cyclicContractJSON)

	assert.PanicsWithValue(
		t,
		"schema Pet is too deep with more than 10 levels",
		func() { dslContract.HydrateContract(contract) },
	)
}
