package model_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/model"
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

func TestProviderHash_BridgesConsumerAndProvider(t *testing.T) {
	providerContract := &model.Contract{Name: "pets-service", Owner: "pets-team"}
	providerKey := providerContract.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", nil))
	provider := providerContract.Resources[providerKey]

	consumerContract := &model.Contract{Name: "shop", Owner: "shop-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil))
	consumer := consumerContract.Resources[consumerKey]

	assert.Equal(t, provider.ProviderHash(), consumer.ProviderHash(),
		"consumer's ProviderHash must equal the provider's ProviderHash for the same endpoint")
}

func TestProviderHash_StatusCodeChangesHash(t *testing.T) {
	c := &model.Contract{Name: "pets-service", Owner: "pets-team"}
	c.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", nil))
	c.AddResource(model.NewProvidedRestResponse("/pets", "get", "404", nil))

	r200 := model.NewProvidedRestResponse("/pets", "get", "200", nil)
	r200.ContractInfo = &model.ContractInfo{Name: "pets-service", Owner: "pets-team"}
	r404 := model.NewProvidedRestResponse("/pets", "get", "404", nil)
	r404.ContractInfo = &model.ContractInfo{Name: "pets-service", Owner: "pets-team"}

	assert.NotEqual(t, r200.ProviderHash(), r404.ProviderHash())
}

func TestProviderHash_RequestAndResponseDiffer(t *testing.T) {
	req := model.NewProvidedRestRequest("/pets", "get", nil)
	req.ContractInfo = &model.ContractInfo{Name: "pets-service", Owner: "pets-team"}
	resp := model.NewProvidedRestResponse("/pets", "get", "200", nil)
	resp.ContractInfo = &model.ContractInfo{Name: "pets-service", Owner: "pets-team"}

	assert.NotEqual(t, req.ProviderHash(), resp.ProviderHash())
}

func TestConsumerHash_IncludesProvider(t *testing.T) {
	a := model.NewConsumedRestResponse("pets-service", "/users", "get", "200", nil)
	a.ContractInfo = &model.ContractInfo{Name: "shop", Owner: "shop-team"}
	b := model.NewConsumedRestResponse("auth-service", "/users", "get", "200", nil)
	b.ContractInfo = &model.ContractInfo{Name: "shop", Owner: "shop-team"}

	assert.NotEqual(t, a.ConsumerHash(), b.ConsumerHash(),
		"ConsumerHash must include Provider so the same consumer can reference the same endpoint from two providers")
}

func TestConsumerHash_EmptyForProvides(t *testing.T) {
	c := &model.Contract{Name: "pets-service", Owner: "pets-team"}
	key := c.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", nil))
	provider := c.Resources[key]

	assert.Empty(t, provider.ConsumerHash())
}

func TestConsumerHash_IncludesContractName(t *testing.T) {
	shopKey := (&model.Contract{Name: "shop", Owner: "shop-team"}).AddResource(model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil))
	otherKey := (&model.Contract{Name: "warehouse", Owner: "ops-team"}).AddResource(model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil))

	assert.NotEqual(t, shopKey, otherKey,
		"ConsumerHash must be unique per consumer service")
}
