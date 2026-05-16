package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewProviderRestRequest(t *testing.T) {
	schema := model.NewSchema()
	args := model.ProviderRestRequestArgs{
		Owner:    "owner",
		Endpoint: "/items",
		Method:   "post",
	}

	got := model.NewProviderRestRequest(args, schema)

	expected := model.ProviderRequest{
		Owner:        "owner",
		Schema:       schema,
		ResourceType: model.ProvidesRestRequest,
		RestRequest:  model.RestRequest{Endpoint: "/items", Method: "post"},
	}
	assert.Equal(t, expected, got)
}
