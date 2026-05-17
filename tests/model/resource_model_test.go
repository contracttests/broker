package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumedRestRequest(t *testing.T) {
	properties := map[string]model.Property{}

	got := model.NewConsumedRestRequest("ledger", "/transactions", "post", properties)

	expected := model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestRequest,
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "post",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
}

func TestNewProvidedRestRequest(t *testing.T) {
	properties := map[string]model.Property{}

	got := model.NewProvidedRestRequest("/items", "post", properties)

	expected := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestRequest,
		Endpoint:   "/items",
		Method:     "post",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
	assert.Empty(t, got.Provider)
	assert.Empty(t, got.StatusCode)
}

func TestNewConsumedRestResponse(t *testing.T) {
	properties := map[string]model.Property{}

	got := model.NewConsumedRestResponse("ledger", "/transactions", "post", "200", properties)

	expected := model.Resource{
		Direction:  model.Consumes,
		Kind:       model.RestResponse,
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "post",
		StatusCode: "200",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
}

func TestNewProvidedRestResponse(t *testing.T) {
	properties := map[string]model.Property{}

	got := model.NewProvidedRestResponse("/items", "get", "200", properties)

	expected := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
	assert.Empty(t, got.Provider)
}
