package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerRestResponse(t *testing.T) {
	schema := model.NewSchema()
	args := model.ConsumerRestResponseArgs{
		Owner:      "owner",
		Provider:   "provider",
		Endpoint:   "/charge",
		Method:     "post",
		StatusCode: "200",
	}

	got := model.NewConsumerRestResponse(args, schema)

	expected := model.ConsumerResponse{
		Owner:        "owner",
		Provider:     "provider",
		Schema:       schema,
		ResourceType: model.ConsumesRestResponse,
		RestResponse: model.RestResponse{Endpoint: "/charge", Method: "post", StatusCode: "200"},
	}
	assert.Equal(t, expected, got)
}
