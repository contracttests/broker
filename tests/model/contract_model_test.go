package model_test

import (
	"encoding/json"
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestContractChecksumIsStable(t *testing.T) {
	c := model.Contract{Name: "items-contract", Owner: "app"}

	first := c.Checksum()
	second := c.Checksum()

	assert.Len(t, first, 64)
	assert.Equal(t, first, second)
}

func TestContractChecksumChangesWithContent(t *testing.T) {
	a := model.Contract{Name: "items-contract", Owner: "app"}
	b := model.Contract{Name: "items-contract", Owner: "billing"}

	assert.NotEqual(t, a.Checksum(), b.Checksum())
}

func TestContractChecksumIsInvariantToResourceInsertionOrder(t *testing.T) {
	properties := map[string]model.Property{
		"root":      {Path: "root", Type: "object"},
		"root.name": {Path: "root.name", Type: "string"},
	}
	r1 := model.NewProvidedRestResponse("/items", "get", "200", properties)
	r2 := model.NewConsumedRestRequest("billing", "/invoices", "post", properties)

	a := model.Contract{Name: "items-contract", Owner: "app"}
	a.AddResource(r1)
	a.AddResource(r2)

	b := model.Contract{Name: "items-contract", Owner: "app"}
	b.AddResource(r2)
	b.AddResource(r1)

	assert.Equal(t, a.Checksum(), b.Checksum())
}

func TestContractChecksumIsInvariantToInputFieldOrder(t *testing.T) {
	nameFirst := `{"name":"items-contract","owner":"app"}`
	ownerFirst := `{"owner":"app","name":"items-contract"}`

	var a, b dsl.Contract
	assert.NoError(t, json.Unmarshal([]byte(nameFirst), &a))
	assert.NoError(t, json.Unmarshal([]byte(ownerFirst), &b))

	contractA := a.ToContractModel()
	contractB := b.ToContractModel()

	assert.Equal(t, contractA.Checksum(), contractB.Checksum())
}

func TestContractChecksumIsInvariantToPropertyOrder(t *testing.T) {
	r1 := model.NewProvidedRestResponse("/items", "get", "200", map[string]model.Property{
		"root":      {Path: "root", Type: "object"},
		"root.name": {Path: "root.name", Type: "string"},
		"root.id":   {Path: "root.id", Type: "string"},
	})
	r2 := model.NewProvidedRestResponse("/items", "get", "200", map[string]model.Property{
		"root.id":   {Path: "root.id", Type: "string"},
		"root":      {Path: "root", Type: "object"},
		"root.name": {Path: "root.name", Type: "string"},
	})

	a := model.Contract{Name: "items-contract", Owner: "app"}
	a.AddResource(r1)

	b := model.Contract{Name: "items-contract", Owner: "app"}
	b.AddResource(r2)

	assert.Equal(t, a.Checksum(), b.Checksum())
}
