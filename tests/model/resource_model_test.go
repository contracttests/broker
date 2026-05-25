package model_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestResource_ProviderHash_BridgesConsumerAndProvider(t *testing.T) {
	consumer := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil)
	provider := model.NewProvidedRestResponse("/pets", "get", "200", nil)
	provider.AddParticipant(model.NewParticipant("pets-service"))

	assert.Equal(t, consumer.ProviderHash(), provider.ProviderHash())
}

func TestResource_ConsumerHash_EmptyForProvidedResource(t *testing.T) {
	provider := model.NewProvidedRestResponse("/pets", "get", "200", nil)
	provider.AddParticipant(model.NewParticipant("pets-service"))

	assert.Empty(t, provider.ConsumerHash())
}

func TestResource_PrimaryHash_ConsumerDirection_EqualsConsumerHash(t *testing.T) {
	consumer := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil)
	consumer.AddParticipant(model.NewParticipant("web-app"))

	assert.Equal(t, consumer.ConsumerHash(), consumer.PrimaryHash())
}

func TestResource_ConsumerHash_RestResponse_IncludesStatusCode(t *testing.T) {
	a := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil)
	a.AddParticipant(model.NewParticipant("web-app"))
	b := model.NewConsumedRestResponse("pets-service", "/pets", "get", "404", nil)
	b.AddParticipant(model.NewParticipant("web-app"))

	assert.NotEqual(t, a.ConsumerHash(), b.ConsumerHash())
}
