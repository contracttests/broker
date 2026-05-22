package model_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttesting/broker/server/internal/features/upload_contract/wireout"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompatibility_MissingTwoProperties_AccumulatesBothBreaks(t *testing.T) {
	consumerContract := &model.Contract{Name: "broken-app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", map[string]model.Property{
		"root":             model.NewProperty("root", "array", false),
		"root[]":           model.NewProperty("root[]", "object", false),
		"root[].uuid":      model.NewProperty("root[].uuid", "string", false),
		"root[].name":      model.NewProperty("root[].name", "string", false),
		"root[].deletedAt": model.NewProperty("root[].deletedAt", "string", false),
		"root[].soldAt":    model.NewProperty("root[].soldAt", "string", false),
	}))
	consumer := consumerContract.Resources[consumerKey]
	provider := model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":        model.NewProperty("root", "array", false),
		"root[]":      model.NewProperty("root[]", "object", false),
		"root[].uuid": model.NewProperty("root[].uuid", "string", false),
		"root[].name": model.NewProperty("root[].name", "string", false),
	})

	got := model.Compare(consumer, provider)

	require.Len(t, got, 2, "Compare must accumulate every missing property; short-circuit is a bug")
	properties := map[string]model.BreakingReason{}
	for _, breakingChange := range got {
		properties[breakingChange.Property] = breakingChange.Reason
		assert.Equal(t, "broken-app", breakingChange.ContractInfo.Name)
		assert.Equal(t, "app-team", breakingChange.ContractInfo.Owner)
	}
	assert.Equal(t, model.ReasonMissingInProvider, properties["root[].deletedAt"])
	assert.Equal(t, model.ReasonMissingInProvider, properties["root[].soldAt"])
}

func TestCompatibility_TypeMismatch_CarriesExpectedAndActualTypes(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestRequest("pets-service", "/pets", "post", map[string]model.Property{
		"root":     model.NewProperty("root", "object", false),
		"root.age": model.NewProperty("root.age", "string", false),
	}))
	consumer := consumerContract.Resources[consumerKey]
	provider := model.NewProvidedRestRequest("/pets", "post", map[string]model.Property{
		"root":     model.NewProperty("root", "object", false),
		"root.age": model.NewProperty("root.age", "integer", false),
	})

	got := model.Compare(consumer, provider)

	require.Len(t, got, 1)
	assert.Equal(t, model.ReasonTypeMismatch, got[0].Reason)
	assert.Equal(t, "root.age", got[0].Property)
	assert.Equal(t, "string", got[0].ExpectedType)
	assert.Equal(t, "integer", got[0].ActualType)
	assert.Equal(t, "app", got[0].ContractInfo.Name)
	assert.Equal(t, "app-team", got[0].ContractInfo.Owner)
}

func TestCompatibility_OptionalInProviderRequiredInConsumer(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
	}))
	consumer := consumerContract.Resources[consumerKey]
	provider := model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", true),
	})

	got := model.Compare(consumer, provider)

	require.Len(t, got, 1)
	assert.Equal(t, model.ReasonOptionalInProviderRequiredInConsumer, got[0].Reason)
	assert.Equal(t, "root.uuid", got[0].Property)
	assert.Equal(t, "app", got[0].ContractInfo.Name)
	assert.Equal(t, "app-team", got[0].ContractInfo.Owner)
}

func TestCompatibility_SubsetConsumer_NoBreaks(t *testing.T) {
	consumer := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
	})
	provider := model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
		"root.name": model.NewProperty("root.name", "string", false),
	})

	got := model.Compare(consumer, provider)
	assert.Empty(t, got)
}

func TestCompatibility_ProviderHasExtras_NoBreaks(t *testing.T) {
	consumer := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
	})
	provider := model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
		"root.name": model.NewProperty("root.name", "string", true),
		"root.age":  model.NewProperty("root.age", "integer", false),
	})

	got := model.Compare(consumer, provider)
	assert.Empty(t, got)
}

func TestCompatibility_ConsumerOptional_ProviderRequired_NoBreak(t *testing.T) {
	consumer := model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", true),
	})
	provider := model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.uuid": model.NewProperty("root.uuid", "string", false),
	})

	got := model.Compare(consumer, provider)
	assert.Empty(t, got)
}

func TestCompatibility_MissingProviderBreak_HasEmptyPropertyAndCopiesCoordsAndConsumerInfo(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "billing-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestResponse("billing", "/invoices", "get", "200", nil))
	consumer := consumerContract.Resources[consumerKey]

	for _, reason := range []model.BreakingReason{
		model.ReasonProviderNotSpecified,
		model.ReasonProviderNotFound,
		model.ReasonProviderResourceNotFound,
	} {
		breakingChange := model.NewMissingProviderBreak(consumer, reason)
		assert.Equal(t, "", breakingChange.Property, "reason=%s", reason)
		assert.Equal(t, reason, breakingChange.Reason)
		assert.Equal(t, model.NewBrokenResource(consumer), breakingChange.Resource)
		assert.Equal(t, "app", breakingChange.ContractInfo.Name, "reason=%s", reason)
		assert.Equal(t, "billing-team", breakingChange.ContractInfo.Owner, "reason=%s", reason)
	}
}

func TestCompatibility_MissingProviderBreak_SerializesWithoutPropertyKey(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestResponse("billing", "/invoices", "get", "200", nil))
	consumer := consumerContract.Resources[consumerKey]
	breakingChange := model.NewMissingProviderBreak(consumer, model.ReasonProviderNotFound)

	view := wireout.BreakingChangeItem{
		ContractName:  "app",
		ContractOwner: "app-team",
		Resource: wireout.BrokenResource{
			Direction:  string(breakingChange.Resource.Direction),
			Kind:       string(breakingChange.Resource.Kind),
			Provider:   breakingChange.Resource.Provider,
			Endpoint:   breakingChange.Resource.Endpoint,
			Method:     breakingChange.Resource.Method,
			StatusCode: breakingChange.Resource.StatusCode,
		},
		Property: breakingChange.Property,
		Reason:   string(breakingChange.Reason),
	}
	encoded, err := json.Marshal(view)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), `"property"`)
}

func TestCompatibility_Request_ProviderRequiredMissingInConsumer_Breaks(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestRequest("pets-service", "/pets", "post", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}))
	consumer := consumerContract.Resources[consumerKey]
	provider := model.NewProvidedRestRequest("/pets", "post", map[string]model.Property{
		"root":       model.NewProperty("root", "object", false),
		"root.name":  model.NewProperty("root.name", "string", false),
		"root.email": model.NewProperty("root.email", "string", false),
	})

	got := model.Compare(consumer, provider)

	require.Len(t, got, 1)
	assert.Equal(t, model.ReasonMissingInConsumer, got[0].Reason)
	assert.Equal(t, "root.email", got[0].Property)
	assert.Equal(t, "app", got[0].ContractInfo.Name)
	assert.Equal(t, "app-team", got[0].ContractInfo.Owner)
}

func TestCompatibility_Request_ProviderOptionalMissingInConsumer_NoBreak(t *testing.T) {
	consumer := model.NewConsumedRestRequest("pets-service", "/pets", "post", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	})
	provider := model.NewProvidedRestRequest("/pets", "post", map[string]model.Property{
		"root":       model.NewProperty("root", "object", false),
		"root.name":  model.NewProperty("root.name", "string", false),
		"root.email": model.NewProperty("root.email", "string", true),
	})

	got := model.Compare(consumer, provider)
	assert.Empty(t, got)
}

func TestCompatibility_Request_ProviderRequiredConsumerOptional_Breaks(t *testing.T) {
	consumerContract := &model.Contract{Name: "app", Owner: "app-team"}
	consumerKey := consumerContract.AddResource(model.NewConsumedRestRequest("pets-service", "/pets", "post", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", true),
	}))
	consumer := consumerContract.Resources[consumerKey]
	provider := model.NewProvidedRestRequest("/pets", "post", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	})

	got := model.Compare(consumer, provider)

	require.Len(t, got, 1)
	assert.Equal(t, model.ReasonOptionalInConsumerRequiredInProvider, got[0].Reason)
	assert.Equal(t, "root.name", got[0].Property)
}

func TestCompatibility_Request_ConsumerExtraField_NoBreak(t *testing.T) {
	consumer := model.NewConsumedRestRequest("pets-service", "/pets", "post", map[string]model.Property{
		"root":       model.NewProperty("root", "object", false),
		"root.name":  model.NewProperty("root.name", "string", false),
		"root.extra": model.NewProperty("root.extra", "string", false),
	})
	provider := model.NewProvidedRestRequest("/pets", "post", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	})

	got := model.Compare(consumer, provider)
	assert.Empty(t, got, "extra fields the provider does not expect must not break a request")
}

func TestContract_Checksum_ExcludesContractInfo(t *testing.T) {
	build := func() model.Contract {
		c := model.Contract{Name: "billing", Owner: "billing-team"}
		c.AddResource(model.NewProvidedRestResponse("/invoices", "get", "200", map[string]model.Property{
			"root":    model.NewProperty("root", "object", false),
			"root.id": model.NewProperty("root.id", "string", false),
		}))
		return c
	}
	first := build()
	second := build()
	assert.Equal(t, first.Checksum(), second.Checksum(), "ContractInfo on resources must not enter the checksum")
}
