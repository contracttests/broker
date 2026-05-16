package dsl_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestToContractModelEmpty(t *testing.T) {
	got := (&dsl.Contract{}).ToContractModel()

	assert.Empty(t, got.ConsumerRequests)
	assert.Empty(t, got.ConsumerResponses)
	assert.Empty(t, got.ProviderRequests)
	assert.Empty(t, got.ProviderResponses)
}

func TestMappingMinimumConsumer(t *testing.T) {
	payload := `
	{
		"api": {
			"name": "consumer"
		},
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

	expectedRequestSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":        {Path: "root", Type: "object", Optional: false},
			"root.amount": {Path: "root.amount", Type: "integer", Optional: false},
		},
	}
	expectedResponseSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":    {Path: "root", Type: "object", Optional: false},
			"root.ok": {Path: "root.ok", Type: "boolean", Optional: false},
		},
	}

	expectedConsumerRequests := []model.ConsumerRequest{
		{
			Owner:        "consumer",
			Provider:     "payments",
			Schema:       expectedRequestSchema,
			ResourceType: model.ConsumesRestRequest,
			RestRequest:  model.RestRequest{Endpoint: "/charge", Method: "post"},
		},
	}

	expectedConsumerResponses := []model.ConsumerResponse{
		{
			Owner:        "consumer",
			Provider:     "payments",
			Schema:       expectedResponseSchema,
			ResourceType: model.ConsumesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/charge", Method: "post", StatusCode: "200"},
		},
	}

	contractModel := contract.ToContractModel()

	assert.ElementsMatch(t, expectedConsumerRequests, contractModel.ConsumerRequests)
	assert.ElementsMatch(t, expectedConsumerResponses, contractModel.ConsumerResponses)
	assert.Empty(t, contractModel.ProviderRequests)
	assert.Empty(t, contractModel.ProviderResponses)
}

func TestMappingMinimumProvider(t *testing.T) {
	payload := `
	{
		"api": {
			"name": "provider"
		},
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

	expectedResponseSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":        {
				Path: "root", 
				Type: "object", 
				Optional: false,
			},
			"root.status": {
				Path: "root.status", 
				Type: "string", 
				Optional: false,
			},
		},
	}

	expectedProviderResponses := []model.ProviderResponse{
		{
			Owner:        "provider",
			Schema:       expectedResponseSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{
				Endpoint:   "/health",
				Method:     "get",
				StatusCode: "200",
			},
		},
	}

	assert.Empty(t, got.ConsumerRequests)
	assert.Empty(t, got.ConsumerResponses)
	assert.Empty(t, got.ProviderRequests)
	assert.ElementsMatch(t, expectedProviderResponses, got.ProviderResponses)
}

func TestMappingFull(t *testing.T) {
	payload := `
	{
		"api": {
			"name": "payments"
		},
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

	paymentRequestSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":             {Path: "root", Type: "object", Optional: false},
			"root.amount":      {Path: "root.amount", Type: "integer", Optional: false},
			"root.currency":    {Path: "root.currency", Type: "string", Optional: true},
			"root.customer":    {Path: "root.customer", Type: "object", Optional: false},
			"root.customer.id": {Path: "root.customer.id", Type: "string", Optional: false},
		},
	}
	paymentSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":        {Path: "root", Type: "object", Optional: false},
			"root.id":     {Path: "root.id", Type: "string", Optional: false},
			"root.status": {Path: "root.status", Type: "string", Optional: true},
		},
	}
	errorSchema := model.Schema{
		Properties: model.SchemaProperties{
			"root":         {Path: "root", Type: "object", Optional: false},
			"root.code":    {Path: "root.code", Type: "string", Optional: false},
			"root.message": {Path: "root.message", Type: "string", Optional: true},
		},
	}

	expectedProviderRequests := []model.ProviderRequest{
		{
			Owner:        "payments",
			Schema:       paymentRequestSchema,
			ResourceType: model.ProvidesRestRequest,
			RestRequest:  model.RestRequest{Endpoint: "/payments", Method: "post"},
		},
		{
			Owner:        "payments",
			Schema:       paymentRequestSchema,
			ResourceType: model.ProvidesRestRequest,
			RestRequest:  model.RestRequest{Endpoint: "/payments", Method: "put"},
		},
	}

	expectedProviderResponses := []model.ProviderResponse{
		{
			Owner:        "payments",
			Schema:       paymentSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "get", StatusCode: "200"},
		},
		{
			Owner:        "payments",
			Schema:       paymentSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "post", StatusCode: "201"},
		},
		{
			Owner:        "payments",
			Schema:       errorSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "post", StatusCode: "400"},
		},
		{
			Owner:        "payments",
			Schema:       paymentSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "put", StatusCode: "200"},
		},
		{
			Owner:        "payments",
			Schema:       errorSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "put", StatusCode: "404"},
		},
		{
			Owner:        "payments",
			Schema:       paymentSchema,
			ResourceType: model.ProvidesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/payments", Method: "delete", StatusCode: "204"},
		},
	}

	expectedConsumerRequests := []model.ConsumerRequest{
		{
			Owner:        "payments",
			Provider:     "ledger",
			Schema:       paymentRequestSchema,
			ResourceType: model.ConsumesRestRequest,
			RestRequest:  model.RestRequest{Endpoint: "/transactions", Method: "post"},
		},
	}

	expectedConsumerResponses := []model.ConsumerResponse{
		{
			Owner:        "payments",
			Provider:     "ledger",
			Schema:       paymentSchema,
			ResourceType: model.ConsumesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/transactions", Method: "get", StatusCode: "200"},
		},
		{
			Owner:        "payments",
			Provider:     "ledger",
			Schema:       paymentSchema,
			ResourceType: model.ConsumesRestResponse,
			RestResponse: model.RestResponse{Endpoint: "/transactions", Method: "post", StatusCode: "202"},
		},
	}

	assert.ElementsMatch(t, expectedProviderRequests, contractModel.ProviderRequests)
	assert.ElementsMatch(t, expectedProviderResponses, contractModel.ProviderResponses)
	assert.ElementsMatch(t, expectedConsumerRequests, contractModel.ConsumerRequests)
	assert.ElementsMatch(t, expectedConsumerResponses, contractModel.ConsumerResponses)
}