package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerRestRequest(t *testing.T) {
	schema := model.NewSchema()

	expected := model.ConsumerRequest{
		Owner:        "payments",
		Provider:     "ledger",
		Schema:       schema,
		ResourceType: model.ConsumesRestRequest,
		RestRequest:  model.RestRequest{Endpoint: "/transactions", Method: "post"},
	}

	args := model.ConsumerRestRequestArgs{
		Owner:    "payments",
		Provider: "ledger",
		Endpoint: "/transactions",
		Method:   "post",
	}

	actual := model.NewConsumerRestRequest(args, schema)

	assert.Equal(t, expected, actual)
}
