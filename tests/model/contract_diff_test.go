package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDiff_ResourceAdded(t *testing.T) {
	prev := model.Contract{Name: "c", Owner: "app"}

	added := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":      {Path: "root", Type: "object"},
			"root.name": {Path: "root.name", Type: "string"},
		},
	}
	next := model.Contract{Name: "c", Owner: "app"}
	key := next.AddResource(added)

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeAdded, got.Kind)
		assert.Len(t, got.Properties, 2)
		for _, p := range got.Properties {
			assert.Equal(t, model.ChangeAdded, p.Kind)
		}
	}
}

func TestDiff_ResourceRemoved(t *testing.T) {
	removed := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":      {Path: "root", Type: "object"},
			"root.name": {Path: "root.name", Type: "string"},
		},
	}
	prev := model.Contract{Name: "c", Owner: "app"}
	key := prev.AddResource(removed)

	next := model.Contract{Name: "c", Owner: "app"}

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeRemoved, got.Kind)
		assert.Len(t, got.Properties, 2)
		for _, p := range got.Properties {
			assert.Equal(t, model.ChangeRemoved, p.Kind)
		}
	}
}

func TestDiff_ResourceChanged(t *testing.T) {
	prev := model.Contract{Name: "c", Owner: "app"}
	prev.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":    {Path: "root", Type: "object"},
			"root.id": {Path: "root.id", Type: "string"},
		},
	})

	changed := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":    {Path: "root", Type: "object"},
			"root.id": {Path: "root.id", Type: "integer"},
		},
	}
	next := model.Contract{Name: "c", Owner: "app"}
	key := next.AddResource(changed)

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeModified, got.Kind)
	}
}

func TestDiff_PropertyAdded(t *testing.T) {
	prev := model.Contract{Name: "c", Owner: "app"}
	prev.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root": {Path: "root", Type: "object"},
		},
	})

	updated := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":      {Path: "root", Type: "object"},
			"root.tags": {Path: "root.tags", Type: "array"},
		},
	}
	next := model.Contract{Name: "c", Owner: "app"}
	key := next.AddResource(updated)

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeModified, got.Kind)
		if assert.Len(t, got.Properties, 1) {
			p := got.Properties["root.tags"]
			assert.Equal(t, model.ChangeAdded, p.Kind)
			assert.Equal(t, "array", p.After.Type)
		}
	}
}

func TestDiff_PropertyRemoved(t *testing.T) {
	prev := model.Contract{Name: "c", Owner: "app"}
	prev.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root":      {Path: "root", Type: "object"},
			"root.tags": {Path: "root.tags", Type: "array"},
		},
	})

	updated := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root": {Path: "root", Type: "object"},
		},
	}
	next := model.Contract{Name: "c", Owner: "app"}
	key := next.AddResource(updated)

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeModified, got.Kind)
		if assert.Len(t, got.Properties, 1) {
			p := got.Properties["root.tags"]
			assert.Equal(t, model.ChangeRemoved, p.Kind)
			assert.Equal(t, "array", p.Before.Type)
		}
	}
}

func TestDiff_PropertyChanged(t *testing.T) {
	prev := model.Contract{Name: "c", Owner: "app"}
	prev.AddResource(model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root.id": {Path: "root.id", Type: "string"},
		},
	})

	updated := model.Resource{
		Direction:  model.Provides,
		Kind:       model.RestResponse,
		Endpoint:   "/items",
		Method:     "get",
		StatusCode: "200",
		Properties: map[string]model.Property{
			"root.id": {Path: "root.id", Type: "integer"},
		},
	}
	next := model.Contract{Name: "c", Owner: "app"}
	key := next.AddResource(updated)

	diff := prev.Diff(&next)

	got, ok := diff.Resources[key]
	if assert.True(t, ok) {
		assert.Equal(t, model.ChangeModified, got.Kind)
		if assert.Len(t, got.Properties, 1) {
			p := got.Properties["root.id"]
			assert.Equal(t, model.ChangeModified, p.Kind)
			assert.Equal(t, "string", p.Before.Type)
			assert.Equal(t, "integer", p.After.Type)
		}
	}
}
