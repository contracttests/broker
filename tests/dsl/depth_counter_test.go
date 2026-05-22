package dsl_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestDepthCounterEnterAtLimitPanics(t *testing.T) {
	c := dsl.NewDepthCounter("anyschema")

	for i := 1; i < dsl.MAX_DEPTH; i++ {
		c.Enter()
	}

	assert.PanicsWithValue(t, "schema anyschema is too deep with more than 10 levels", func() {
		c.Enter()
	})
}

func TestToContractModelPanicsWhenSchemaIsTooDeep(t *testing.T) {
	payload := `
	{
		"owner": "test",
		"provides": {
			"rest": {
				"/deep": {
					"get": {
						"responses": {
							"200": "L0"
						}
					}
				}
			}
		},
		"schemas": {
			"L0": {
				"type": "object",
				"properties": {
					"a": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"b": {
									"$ref": "L1"
								}
							}
						}
					}
				}
			},
			"L1": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"c": {
							"$ref": "L2"
						}
					}
				}
			},
			"L2": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"d": {
							"type": "object",
							"properties": {
								"e": {
									"type": "string"
								}
							}
						}
					}
				}
			}
		}
	}`

	var contract dsl.Contract
	if !assert.NoError(t, json.Unmarshal([]byte(payload), &contract)) {
		return
	}

	assert.PanicsWithValue(t, "schema L0 is too deep with more than 10 levels", func() {
		contract.ToContractModel()
	})
}
