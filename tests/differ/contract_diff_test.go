package differ_test

import (
	"testing"

	"github.com/contracttesting/broker/internal/contract_differ"
	"github.com/contracttesting/broker/internal/model"
	"github.com/stretchr/testify/assert"
)

func newContractWithOnePetsResource(participantName string) *model.Contract {
	contract := model.NewContract(model.NewParticipant(participantName), "1", "raw")
	contract.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}))
	return contract
}

func TestDiff_NoChanges_BetweenEquivalentContracts(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")
	next := newContractWithOnePetsResource("pets-service")

	diff := contract_differ.Diff(prev, next)

	assert.Empty(t, diff.Resources)
}

func TestDiff_ReportsAddedResource(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")

	next := newContractWithOnePetsResource("pets-service")
	added := model.NewProvidedRestResponse("/pets/{id}", "get", "200", map[string]model.Property{
		"root": model.NewProperty("root", "object", false),
	})
	next.AddResource(added)

	diff := contract_differ.Diff(prev, next)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeAdded, change.Kind)
		assert.Equal(t, "/pets/{id}", change.Resource.Endpoint)
	}
}

func TestDiff_NextNil_AllResourcesRemoved(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")

	diff := contract_differ.Diff(prev, nil)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeRemoved, change.Kind)
	}
}

func TestDiff_RemovedResource(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")
	extra := model.NewProvidedRestResponse("/pets/{id}", "get", "200", map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	})
	prev.AddResource(extra)

	next := newContractWithOnePetsResource("pets-service")

	diff := contract_differ.Diff(prev, next)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeRemoved, change.Kind)
		assert.Equal(t, "/pets/{id}", change.Resource.Endpoint)
		assert.Len(t, change.Properties, 2)
		for _, propChange := range change.Properties {
			assert.Equal(t, contract_differ.ChangeRemoved, propChange.Kind)
		}
	}
}

func TestDiff_ModifiedResource_PropertyAdded(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")

	next := model.NewContract(model.NewParticipant("pets-service"), "1", "raw")
	next.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.id":   model.NewProperty("root.id", "string", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}))

	diff := contract_differ.Diff(prev, next)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeModified, change.Kind)
		added, ok := change.Properties["root.name"]
		if assert.True(t, ok, "expected root.name in property changes") {
			assert.Equal(t, contract_differ.ChangeAdded, added.Kind)
			assert.Equal(t, "root.name", added.After.Path)
		}
	}
}

func TestDiff_ModifiedResource_PropertyRemoved(t *testing.T) {
	prev := model.NewContract(model.NewParticipant("pets-service"), "1", "raw")
	prev.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":      model.NewProperty("root", "object", false),
		"root.id":   model.NewProperty("root.id", "string", false),
		"root.name": model.NewProperty("root.name", "string", false),
	}))

	next := newContractWithOnePetsResource("pets-service")

	diff := contract_differ.Diff(prev, next)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeModified, change.Kind)
		removed, ok := change.Properties["root.name"]
		if assert.True(t, ok, "expected root.name in property changes") {
			assert.Equal(t, contract_differ.ChangeRemoved, removed.Kind)
			assert.Equal(t, "root.name", removed.Before.Path)
		}
	}
}

func TestDiff_ModifiedResource_PropertyTypeChanged(t *testing.T) {
	prev := newContractWithOnePetsResource("pets-service")

	next := model.NewContract(model.NewParticipant("pets-service"), "1", "raw")
	next.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "int", false),
	}))

	diff := contract_differ.Diff(prev, next)

	assert.Len(t, diff.Resources, 1)
	for _, change := range diff.Resources {
		assert.Equal(t, contract_differ.ChangeModified, change.Kind)
		modified, ok := change.Properties["root.id"]
		if assert.True(t, ok, "expected root.id in property changes") {
			assert.Equal(t, contract_differ.ChangeModified, modified.Kind)
			assert.Equal(t, "string", modified.Before.Type)
			assert.Equal(t, "int", modified.After.Type)
		}
	}
}
