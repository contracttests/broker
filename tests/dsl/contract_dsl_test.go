package dsl_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestToContractModelEmpty(t *testing.T) {
	got := (&dsl.Contract{}).ToContractModel()

	assert.Empty(t, got.Resources)
}

func TestNamePassthrough(t *testing.T) {
	payload := `{"name": "items-contract", "owner": "app"}`

	var contract dsl.Contract
	if !assert.NoError(t, json.Unmarshal([]byte(payload), &contract)) {
		return
	}

	got := contract.ToContractModel()

	assert.Equal(t, "items-contract", got.Name)
	assert.Equal(t, "app", got.Owner)
}

func TestMappingMinimumConsumer(t *testing.T) {
	payload := `
	{
		"owner": "consumer",
		"consumes": {
			"payments": {
				"rest": {
					"/charge": {
						"post": {
							"request": "ChargeRequest",
							"responses": {
								"200": "ChargeResponse"
							}
						}
					}
				}
			}
		},
		"schemas": {
			"ChargeRequest": {
				"type": "object",
				"properties": {
					"amount": {
						"type": "integer"
					}
				}
			},
			"ChargeResponse": {
				"type": "object",
				"properties": {
					"ok": {
						"type": "boolean"
					}
				}
			}
		}
	}`

	var contract dsl.Contract
	if !assert.NoError(t, json.Unmarshal([]byte(payload), &contract)) {
		return
	}

	expectedRequestProperties := map[string]model.Property{
		"root":        {Path: "root", Type: "object", Optional: false},
		"root.amount": {Path: "root.amount", Type: "integer", Optional: false},
	}
	expectedResponseProperties := map[string]model.Property{
		"root":    {Path: "root", Type: "object", Optional: false},
		"root.ok": {Path: "root.ok", Type: "boolean", Optional: false},
	}

	expected := model.Contract{Owner: "consumer"}
	expected.AddResource(model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestRequest,
		Provider:   "payments",
		Endpoint:   "/charge",
		Method:     "post",
		Properties: expectedRequestProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestResponse,
		Provider:   "payments",
		Endpoint:   "/charge",
		Method:     "post",
		StatusCode: "200",
		Properties: expectedResponseProperties,
	})

	assert.Equal(t, expected, contract.ToContractModel())
}

func TestMappingMinimumProvider(t *testing.T) {
	payload := `
	{
		"owner": "provider",
		"provides": {
			"rest": {
				"/health": {
					"get": {
						"responses": {
							"200": "Health"
						}
					}
				}
			}
		},
		"schemas": {
			"Health": {
				"type": "object",
				"properties": {
					"status": {
						"type": "string"
					}
				}
			}
		}
	}`

	var contract dsl.Contract
	if !assert.NoError(t, json.Unmarshal([]byte(payload), &contract)) {
		return
	}

	got := contract.ToContractModel()

	expectedResponseProperties := map[string]model.Property{
		"root": {
			Path:     "root",
			Type:     "object",
			Optional: false,
		},
		"root.status": {
			Path:     "root.status",
			Type:     "string",
			Optional: false,
		},
	}

	expected := model.Contract{Owner: "provider"}
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/health",
		Method:     "get",
		StatusCode: "200",
		Properties: expectedResponseProperties,
	})

	assert.Equal(t, expected, got)
}

func TestMappingFull(t *testing.T) {
	payload := `
	{
		"owner": "payments",
		"provides": {
			"rest": {
				"/payments": {
					"get": {
						"responses": {
							"200": "Payment"
						}
					},
					"post": {
						"request": "PaymentRequest",
						"responses": {
							"201": "Payment",
							"400": "Error"
						}
					},
					"put": {
						"request": "PaymentRequest",
						"responses": {
							"200": "Payment",
							"404": "Error"
						}
					},
					"delete": {
						"responses": {
							"204": "Payment"
						}
					}
				}
			}
		},
		"consumes": {
			"ledger": {
				"rest": {
					"/transactions": {
						"get": {
							"responses": {
								"200": "Payment"
							}
						},
						"post": {
							"request": "PaymentRequest",
							"responses": {
								"202": "Payment"
							}
						}
					}
				}
			}
		},
		"schemas": {
			"PaymentRequest": {
				"type": "object",
				"properties": {
					"amount": {
						"type": "integer"
					},
					"currency": {
						"type": "string",
						"optional": true
					},
					"customer": {
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
			"Payment": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string"
					},
					"status": {
						"type": "string",
						"optional": true
					}
				}
			},
			"Error": {
				"type": "object",
				"properties": {
					"code": {
						"type": "string"
					},
					"message": {
						"type": "string",
						"optional": true
					}
				}
			}
		}
	}`

	var contract dsl.Contract
	if !assert.NoError(t, json.Unmarshal([]byte(payload), &contract)) {
		return
	}
	contractModel := contract.ToContractModel()

	paymentRequestProperties := map[string]model.Property{
		"root":             {Path: "root", Type: "object", Optional: false},
		"root.amount":      {Path: "root.amount", Type: "integer", Optional: false},
		"root.currency":    {Path: "root.currency", Type: "string", Optional: true},
		"root.customer":    {Path: "root.customer", Type: "object", Optional: false},
		"root.customer.id": {Path: "root.customer.id", Type: "string", Optional: false},
	}
	paymentProperties := map[string]model.Property{
		"root":        {Path: "root", Type: "object", Optional: false},
		"root.id":     {Path: "root.id", Type: "string", Optional: false},
		"root.status": {Path: "root.status", Type: "string", Optional: true},
	}
	errorProperties := map[string]model.Property{
		"root":         {Path: "root", Type: "object", Optional: false},
		"root.code":    {Path: "root.code", Type: "string", Optional: false},
		"root.message": {Path: "root.message", Type: "string", Optional: true},
	}

	expected := model.Contract{Owner: "payments"}
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestRequest,
		Endpoint:   "/payments",
		Method:     "post",
		Properties: paymentRequestProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestRequest,
		Endpoint:   "/payments",
		Method:     "put",
		Properties: paymentRequestProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "get",
		StatusCode: "200",
		Properties: paymentProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "post",
		StatusCode: "201",
		Properties: paymentProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "post",
		StatusCode: "400",
		Properties: errorProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "put",
		StatusCode: "200",
		Properties: paymentProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "put",
		StatusCode: "404",
		Properties: errorProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/payments",
		Method:     "delete",
		StatusCode: "204",
		Properties: paymentProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestRequest,
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "post",
		Properties: paymentRequestProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestResponse,
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "get",
		StatusCode: "200",
		Properties: paymentProperties,
	})
	expected.AddResource(model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestResponse,
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "post",
		StatusCode: "202",
		Properties: paymentProperties,
	})

	assert.Equal(t, expected, contractModel)
}
