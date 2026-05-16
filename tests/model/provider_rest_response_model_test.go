package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewProviderRestResponse(t *testing.T) {
	schema := model.NewSchema()
	args := model.ProviderRestResponseArgs{
		Owner:      "owner",
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
	}

	got := model.NewProviderRestResponse(args, schema)

	expected := model.ProviderResponse{
		Owner:        "owner",
		Schema:       schema,
		ResourceType: model.ProvidesRestResponse,
		RestResponse: model.RestResponse{Endpoint: "/items", Method: "get", StatusCode: "200"},
	}
	assert.Equal(t, expected, got)
}
