package dsl_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/model"
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
	p := dsl.NewResourcePath("app;provides;rest;/items;post:request")
	assert.True(t, p.IsProvider())
	assert.False(t, p.IsConsumer())
}

func TestResourcePathIsProviderResponse(t *testing.T) {
	p := dsl.NewResourcePath("app;provides;rest;/items;post:responses;200")
	assert.True(t, p.IsProvider())
	assert.False(t, p.IsConsumer())
}

func TestResourcePathIsConsumerRequest(t *testing.T) {
	p := dsl.NewResourcePath("app;consumes;api;rest;/items;post:request")
	assert.False(t, p.IsProvider())
	assert.True(t, p.IsConsumer())
}

func TestResourcePathIsConsumerResponse(t *testing.T) {
	p := dsl.NewResourcePath("app;consumes;api;rest;/items;post:responses;200")
	assert.False(t, p.IsProvider())
	assert.True(t, p.IsConsumer())
}

func TestResourcePathToConsumerRestRequestArgs(t *testing.T) {
	p := dsl.NewResourcePath("payments;consumes;ledger;rest;/transactions;post;request")
	args := p.ToConsumerRestRequestArgs()
	assert.Equal(t, model.ConsumerRestRequestArgs{
		Owner:    "payments",
		Provider: "ledger",
		Endpoint: "/transactions",
		Method:   "post",
	}, args)
}

func TestResourcePathToConsumerRestResponseArgs(t *testing.T) {
	p := dsl.NewResourcePath("payments;consumes;ledger;rest;/transactions;post;responses;200")
	args := p.ToConsumerRestResponseArgs()
	assert.Equal(t, model.ConsumerRestResponseArgs{
		Owner:      "payments",
		Provider:   "ledger",
		Endpoint:   "/transactions",
		Method:     "post",
		StatusCode: "200",
	}, args)
}

func TestResourcePathToProviderRestRequestArgs(t *testing.T) {
	p := dsl.NewResourcePath("payments;provides;rest;/transactions;post;request")
	args := p.ToProviderRestRequestArgs()
	assert.Equal(t, model.ProviderRestRequestArgs{
		Owner:    "payments",
		Endpoint: "/transactions",
		Method:   "post",
	}, args)
}

func TestResourcePathToProviderRestResponseArgs(t *testing.T) {
	p := dsl.NewResourcePath("payments;provides;rest;/transactions;post;responses;200")
	args := p.ToProviderRestResponseArgs()
	assert.Equal(t, model.ProviderRestResponseArgs{
		Owner:      "payments",
		Endpoint:   "/transactions",
		Method:     "post",
		StatusCode: "200",
	}, args)
}


func TestResourcePathToConsumerRestRequestArgsInvalid(t *testing.T) {
	p := dsl.NewResourcePath("payments;consumes;ledger;rest;/transactions;post")
	assert.Panics(t, func() {
		p.ToConsumerRestRequestArgs()
	})
}

func TestResourcePathToConsumerRestResponseArgsInvalid(t *testing.T) {
	p := dsl.NewResourcePath("payments;consumes;ledger;rest;/transactions;responses;200")
	assert.Panics(t, func() {
		p.ToConsumerRestResponseArgs()
	})
}

func TestResourcePathToProviderRestRequestArgsInvalid(t *testing.T) {
	p := dsl.NewResourcePath("payments;provide;rest;/transactions;post;request")
	assert.Panics(t, func() {
		p.ToProviderRestRequestArgs()
	})
}

func TestResourcePathToProviderRestResponseArgsInvalid(t *testing.T) {
	p := dsl.NewResourcePath("payments;rest;/transactions;post;responses;200")
	assert.Panics(t, func() {
		p.ToProviderRestResponseArgs()
	})
}