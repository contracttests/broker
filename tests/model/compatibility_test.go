package model_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCompare_Compatible_NoBreaks(t *testing.T) {
	props := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", props)
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", props)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Empty(t, breaks)
}

func TestCompare_TypeMismatch_ReportsBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}
	providerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "int", false),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", consumerProps)
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", providerProps)
	provider.AddParticipant(model.NewParticipant("pets-service"))

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonTypeMismatch, breaks[0].Reason)
	assert.Equal(t, "root.id", breaks[0].Property)
	assert.Equal(t, "string", breaks[0].ConsumerType)
	assert.Equal(t, "int", breaks[0].ProviderType)
	assert.NotNil(t, breaks[0].Counterpart)
	assert.Equal(t, model.UploaderProvider, breaks[0].Counterpart.Role)
	assert.Equal(t, "pets-service", breaks[0].Counterpart.Name)
}

func TestCompare_Response_MissingInProvider_ReportsBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}
	providerProps := map[string]model.Property{
		"root": model.NewProperty("root", "object", false),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", consumerProps)
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonMissingInProvider, breaks[0].Reason)
	assert.Equal(t, "root.name", breaks[0].Property)
}

func TestCompare_Response_OptionalInProviderRequiredInConsumer_ReportsBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}
	providerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", true),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", consumerProps)
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonOptionalInProviderRequiredInConsumer, breaks[0].Reason)
	assert.Equal(t, "root.id", breaks[0].Property)
}

func TestCompare_Request_MissingRequiredProperty_ReportsMissingInConsumer(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root": model.NewProperty("root", "object", false),
	}
	providerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}

	consumer := *model.NewConsumedRestRequest("pets-service", "/pets", "post", consumerProps)
	provider := *model.NewProvidedRestRequest("/pets", "post", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonMissingInConsumer, breaks[0].Reason)
	assert.Equal(t, "root.name", breaks[0].Property)
}

func TestCompare_Request_MissingOptionalProperty_NoBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root": model.NewProperty("root", "object", false),
	}
	providerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.note": model.NewProperty("root.note", "string", true),
	}

	consumer := *model.NewConsumedRestRequest("pets-service", "/pets", "post", consumerProps)
	provider := *model.NewProvidedRestRequest("/pets", "post", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Empty(t, breaks)
}

func TestCompare_Request_TypeMismatch_ReportsBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "int", false),
	}
	providerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}

	consumer := *model.NewConsumedRestRequest("pets-service", "/pets", "post", consumerProps)
	provider := *model.NewProvidedRestRequest("/pets", "post", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonTypeMismatch, breaks[0].Reason)
	assert.Equal(t, "root.name", breaks[0].Property)
	assert.Equal(t, "int", breaks[0].ConsumerType)
	assert.Equal(t, "string", breaks[0].ProviderType)
}

func TestCompare_Request_OptionalInConsumerRequiredInProvider_ReportsBreak(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", true),
	}
	providerProps := map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}

	consumer := *model.NewConsumedRestRequest("pets-service", "/pets", "post", consumerProps)
	provider := *model.NewProvidedRestRequest("/pets", "post", providerProps)

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderConsumer,
	})

	assert.Len(t, breaks, 1)
	assert.Equal(t, model.ReasonOptionalInConsumerRequiredInProvider, breaks[0].Reason)
	assert.Equal(t, "root.name", breaks[0].Property)
}

func TestCompare_UploaderProvider_CounterpartIsConsumer(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}
	providerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "int", false),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", consumerProps)
	consumer.AddParticipant(model.NewParticipant("web-app"))
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", providerProps)
	provider.AddParticipant(model.NewParticipant("pets-service"))

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderProvider,
	})

	assert.Len(t, breaks, 1)
	assert.NotNil(t, breaks[0].Counterpart)
	assert.Equal(t, model.UploaderConsumer, breaks[0].Counterpart.Role)
	assert.Equal(t, "web-app", breaks[0].Counterpart.Name)
}

func TestCompare_NilParticipant_CounterpartIsNil(t *testing.T) {
	consumerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}
	providerProps := map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "int", false),
	}

	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", consumerProps)
	provider := *model.NewProvidedRestResponse("/pets", "get", "200", providerProps)
	provider.AddParticipant(model.NewParticipant("pets-service"))

	breaks := model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: model.UploaderProvider,
	})

	assert.Len(t, breaks, 1)
	assert.Nil(t, breaks[0].Counterpart)
}

func TestNewMissingProviderBreak_ProviderNotSpecified_NoCounterpart(t *testing.T) {
	consumer := *model.NewConsumedRestResponse("", "/pets", "get", "200", nil)

	change := model.NewMissingProviderBreak(consumer, model.ReasonProviderNotSpecified)

	assert.Equal(t, model.ReasonProviderNotSpecified, change.Reason)
	assert.Equal(t, model.UploaderConsumer, change.UploaderRole)
	assert.Nil(t, change.Counterpart)
	assert.Equal(t, model.RestResponse, change.Resource.Kind)
	assert.Equal(t, "/pets", change.Resource.Endpoint)
	assert.Equal(t, "get", change.Resource.Method)
	assert.Equal(t, "200", change.Resource.StatusCode)
}

func TestNewMissingProviderBreak_ProviderNotFound_NameOnlyCounterpart(t *testing.T) {
	consumer := *model.NewConsumedRestResponse("unknown-service", "/pets", "get", "200", nil)

	change := model.NewMissingProviderBreak(consumer, model.ReasonProviderNotFound)

	assert.Equal(t, model.ReasonProviderNotFound, change.Reason)
	assert.Equal(t, model.UploaderConsumer, change.UploaderRole)
	if assert.NotNil(t, change.Counterpart) {
		assert.Equal(t, model.UploaderProvider, change.Counterpart.Role)
		assert.Equal(t, "unknown-service", change.Counterpart.Name)
	}
	assert.Equal(t, "unknown-service", change.Resource.Provider)
}

func TestNewMissingProviderBreak_ProviderResourceNotFound_NameOnlyCounterpart(t *testing.T) {
	consumer := *model.NewConsumedRestResponse("pets-service", "/pets", "get", "200", nil)

	change := model.NewMissingProviderBreak(consumer, model.ReasonProviderResourceNotFound)

	assert.Equal(t, model.ReasonProviderResourceNotFound, change.Reason)
	assert.Equal(t, model.UploaderConsumer, change.UploaderRole)
	if assert.NotNil(t, change.Counterpart) {
		assert.Equal(t, model.UploaderProvider, change.Counterpart.Role)
		assert.Equal(t, "pets-service", change.Counterpart.Name)
	}
}

func TestCompatibilityReport_Breaks_EmptyReportReturnsNonNilEmptySlice(t *testing.T) {
	report := &model.CompatibilityReport{}

	breaks := report.Breaks()

	assert.NotNil(t, breaks)
	assert.Len(t, breaks, 0)
}

func TestCompatibilityReport_Breaks_PreservesAppendOrder(t *testing.T) {
	report := &model.CompatibilityReport{}
	first := model.BreakingChange{Reason: model.ReasonMissingInProvider, Property: "root.a"}
	second := model.BreakingChange{Reason: model.ReasonTypeMismatch, Property: "root.b"}

	report.Append(first)
	report.Append(second)

	breaks := report.Breaks()
	assert.Len(t, breaks, 2)
	assert.Equal(t, first, breaks[0])
	assert.Equal(t, second, breaks[1])
}
