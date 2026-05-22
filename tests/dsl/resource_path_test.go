package dsl_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestResourcePathAppendFromEmpty(t *testing.T) {
	p := dsl.NewResourcePath("")
	p = p.Append("foo", "bar")
	assert.Equal(t, "foo;bar", p.String())
}

func TestResourcePathAppendFromNonEmpty(t *testing.T) {
	p := dsl.NewResourcePath("a")
	p = p.Append("b", "c")
	assert.Equal(t, "a;b;c", p.String())
}

func TestResourcePathSplit(t *testing.T) {
	p := dsl.NewResourcePath("a;b;c")
	assert.Equal(t, []string{"a", "b", "c"}, p.Split())
}

func TestResourcePathIsProviderRequest(t *testing.T) {
	p := dsl.NewResourcePath("provides;rest;/items;post;request")
	assert.True(t, p.IsProvider())
	assert.False(t, p.IsConsumer())
}

func TestResourcePathIsProviderResponse(t *testing.T) {
	p := dsl.NewResourcePath("provides;rest;/items;post;responses;200")
	assert.True(t, p.IsProvider())
	assert.False(t, p.IsConsumer())
}

func TestResourcePathIsConsumerRequest(t *testing.T) {
	p := dsl.NewResourcePath("consumes;billing;rest;/items;post;request")
	assert.False(t, p.IsProvider())
	assert.True(t, p.IsConsumer())
}

func TestResourcePathIsConsumerResponse(t *testing.T) {
	p := dsl.NewResourcePath("consumes;billing;rest;/items;post;responses;200")
	assert.False(t, p.IsProvider())
	assert.True(t, p.IsConsumer())
}

func TestResourcePathToResourceConsumerRequest(t *testing.T) {
	p := dsl.NewResourcePath("consumes;ledger;rest;/transactions;post;request")
	properties := map[string]model.Property{}

	got := p.ToResource(properties)

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

func TestResourcePathToResourceConsumerResponse(t *testing.T) {
	p := dsl.NewResourcePath("consumes;ledger;rest;/transactions;post;responses;200")
	properties := map[string]model.Property{}

	got := p.ToResource(properties)

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

func TestResourcePathToResourceProviderRequest(t *testing.T) {
	p := dsl.NewResourcePath("provides;rest;/transactions;post;request")
	properties := map[string]model.Property{}

	got := p.ToResource(properties)

	expected := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestRequest,
		Endpoint:   "/transactions",
		Method:     "post",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
}

func TestResourcePathToResourceProviderResponse(t *testing.T) {
	p := dsl.NewResourcePath("provides;rest;/transactions;post;responses;200")
	properties := map[string]model.Property{}

	got := p.ToResource(properties)

	expected := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/transactions",
		Method:     "post",
		StatusCode: "200",
		Properties: properties,
	}
	assert.Equal(t, expected, got)
}

func TestResourcePathToResourceUnrecognized(t *testing.T) {
	p := dsl.NewResourcePath("rest;/transactions;post;responses;200")
	assert.Panics(t, func() {
		p.ToResource(map[string]model.Property{})
	})
}
