package dsl_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestSchemaIsObject(t *testing.T) {
	s := dsl.Schema{Type: "object"}
	assert.True(t, s.IsObject())

	s = dsl.Schema{Type: "string"}
	assert.False(t, s.IsObject())

	s = dsl.Schema{Properties: map[string]dsl.Schema{"field": {Type: "string"}}}
	assert.True(t, s.IsObject())

	s = dsl.Schema{}
	assert.False(t, s.IsObject())
}

func TestSchemaIsArray(t *testing.T) {
	s := dsl.Schema{Type: "array"}
	assert.True(t, s.IsArray())

	s = dsl.Schema{Type: "string"}
	assert.False(t, s.IsArray())

	s = dsl.Schema{Items: &dsl.Schema{Type: "object"}}
	assert.True(t, s.IsArray())

	s = dsl.Schema{}
	assert.False(t, s.IsArray())
}

func TestSchemaIsPrimitive(t *testing.T) {
	s := dsl.Schema{Type: "string"}
	assert.True(t, s.IsPrimitive())

	s = dsl.Schema{Type: "integer"}
	assert.True(t, s.IsPrimitive())

	s = dsl.Schema{Type: "float"}
	assert.True(t, s.IsPrimitive())

	s = dsl.Schema{Type: "number"}
	assert.True(t, s.IsPrimitive())

	s = dsl.Schema{Type: "boolean"}
	assert.True(t, s.IsPrimitive())

	s = dsl.Schema{Type: "object"}
	assert.False(t, s.IsPrimitive())

	s = dsl.Schema{}
	assert.False(t, s.IsPrimitive())
}

func TestSchemaIsRef(t *testing.T) {
	s := dsl.Schema{Ref: "User"}
	assert.True(t, s.IsRef())

	s = dsl.Schema{Ref: ""}
	assert.False(t, s.IsRef())

	s = dsl.Schema{Type: "object", Ref: "User"}
	assert.False(t, s.IsRef())

	s = dsl.Schema{Properties: map[string]dsl.Schema{"field": {Type: "string"}}}
	assert.False(t, s.IsRef())

	s = dsl.Schema{Items: &dsl.Schema{}, Ref: "User"}
	assert.False(t, s.IsRef())
}
