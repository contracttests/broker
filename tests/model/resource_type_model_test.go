package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestResourceTypeConstants(t *testing.T) {
	assert.Equal(t, model.ResourceType("consumes_rest_request"), model.ConsumesRestRequest)
	assert.Equal(t, model.ResourceType("provides_rest_request"), model.ProvidesRestRequest)
	assert.Equal(t, model.ResourceType("consumes_rest_response"), model.ConsumesRestResponse)
	assert.Equal(t, model.ResourceType("provides_rest_response"), model.ProvidesRestResponse)
}
