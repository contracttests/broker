package model_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func newContractWithOnePetsResource(participantName string) *model.Contract {
	contract := model.NewContract(model.NewParticipant(participantName), "raw")
	contract.AddResource(model.NewProvidedRestResponse("/pets", "get", "200", map[string]model.Property{
		"root":    model.NewProperty("root", "object", false),
		"root.id": model.NewProperty("root.id", "string", false),
	}))
	return contract
}

func TestContract_Checksum_IsStableForEquivalentContracts(t *testing.T) {
	a := newContractWithOnePetsResource("pets-service")
	b := newContractWithOnePetsResource("pets-service")

	assert.Equal(t, a.Checksum(), b.Checksum())
}

func TestContract_Checksum_DiffersWhenResourceAdded(t *testing.T) {
	a := newContractWithOnePetsResource("pets-service")

	b := newContractWithOnePetsResource("pets-service")
	b.AddResource(model.NewProvidedRestResponse("/pets/{id}", "get", "200", map[string]model.Property{
		"root": model.NewProperty("root", "object", false),
	}))

	assert.NotEqual(t, a.Checksum(), b.Checksum())
}
