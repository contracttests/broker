package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPropertyCanonicalKey_Required(t *testing.T) {
	p := model.NewProperty("root.id", "integer", false)
	assert.Equal(t, "root.id;;integer;;false", p.CanonicalKey())
}

func TestPropertyCanonicalKey_Optional(t *testing.T) {
	p := model.NewProperty("root.name", "string", true)
	assert.Equal(t, "root.name;;string;;true", p.CanonicalKey())
}

func TestResourceCanonicalKey_ProvidedResponse(t *testing.T) {
	r := model.NewProvidedRestResponse("/items", "get", "200", nil)
	assert.Equal(t, "provides;;rest_response;;;;/items;;get;;200;;", r.CanonicalKey())
}

func TestResourceCanonicalKey_ConsumedRequest(t *testing.T) {
	r := model.NewConsumedRestRequest("billing", "/invoices", "post", nil)
	assert.Equal(t, "consumes;;rest_request;;billing;;/invoices;;post;;;;", r.CanonicalKey())
}

func TestResourceCanonicalKey_WithProperties(t *testing.T) {
	r := model.NewProvidedRestResponse("/items", "get", "200", map[string]model.Property{
		"root":    {Path: "root", Type: "object"},
		"root.id": {Path: "root.id", Type: "string"},
	})
	assert.Equal(t,
		"provides;;rest_response;;;;/items;;get;;200;;root.id;;string;;false;;root;;object;;false",
		r.CanonicalKey(),
	)
}

func TestResourceCanonicalKey_ExcludesContractInfo(t *testing.T) {
	r := model.NewProvidedRestResponse("/items", "get", "200", nil)
	r.ContractInfo = &model.ContractInfo{Name: "service-a", Owner: "team-a"}
	assert.Equal(t, "provides;;rest_response;;;;/items;;get;;200;;", r.CanonicalKey())
}

func TestResourceCanonicalKey_PropertyOrderInvariant(t *testing.T) {
	expected := "provides;;rest_response;;;;/items;;get;;200;;root.id;;string;;false;;root;;object;;false"

	r1 := model.NewProvidedRestResponse("/items", "get", "200", map[string]model.Property{
		"root":    {Path: "root", Type: "object"},
		"root.id": {Path: "root.id", Type: "string"},
	})
	assert.Equal(t, expected, r1.CanonicalKey())

	r2 := model.NewProvidedRestResponse("/items", "get", "200", map[string]model.Property{
		"root.id": {Path: "root.id", Type: "string"},
		"root":    {Path: "root", Type: "object"},
	})
	assert.Equal(t, expected, r2.CanonicalKey())
}

func TestContractCanonicalKey_NoResources(t *testing.T) {
	c := model.Contract{Name: "shop", Owner: "shop-team"}
	assert.Equal(t, "shop;;shop-team;;", c.CanonicalKey())
}

func TestContractCanonicalKey_WithSingleResource(t *testing.T) {
	c := model.Contract{Name: "shop", Owner: "shop-team"}
	c.AddResource(model.NewProvidedRestResponse("/items", "get", "200", nil))

	assert.Equal(t,
		"shop;;shop-team;;provides;;rest_response;;;;/items;;get;;200;;",
		c.CanonicalKey(),
	)
}

func TestContractCanonicalKey_ResourceOrderInvariant(t *testing.T) {
	expected := "shop;;shop-team;;consumes;;rest_request;;billing;;/invoices;;post;;;;;;provides;;rest_response;;;;/items;;get;;200;;"

	a := model.Contract{Name: "shop", Owner: "shop-team"}
	a.AddResource(model.NewProvidedRestResponse("/items", "get", "200", nil))
	a.AddResource(model.NewConsumedRestRequest("billing", "/invoices", "post", nil))
	assert.Equal(t, expected, a.CanonicalKey())

	b := model.Contract{Name: "shop", Owner: "shop-team"}
	b.AddResource(model.NewConsumedRestRequest("billing", "/invoices", "post", nil))
	b.AddResource(model.NewProvidedRestResponse("/items", "get", "200", nil))
	assert.Equal(t, expected, b.CanonicalKey())
}
