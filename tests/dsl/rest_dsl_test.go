package dsl_test

import (
	"testing"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetMethodIsNonZero(t *testing.T) {
	m := dsl.GetMethod{}
	assert.False(t, m.IsNonZero())

	m = dsl.GetMethod{Responses: dsl.Responses{200: "Response"}}
	assert.True(t, m.IsNonZero())
}

func TestPostMethodIsNonZeroAndHasRequestBody(t *testing.T) {
	m := dsl.PostMethod{}
	assert.False(t, m.IsNonZero())
	assert.False(t, m.HasRequestBody())

	m = dsl.PostMethod{RequestBody: "Request"}
	assert.True(t, m.IsNonZero())
	assert.True(t, m.HasRequestBody())

	m = dsl.PostMethod{Responses: dsl.Responses{200: "Response"}}
	assert.True(t, m.IsNonZero())
	assert.False(t, m.HasRequestBody())

	m = dsl.PostMethod{RequestBody: "Request", Responses: dsl.Responses{200: "Response"}}
	assert.True(t, m.IsNonZero())
	assert.True(t, m.HasRequestBody())
}

func TestPutMethodIsNonZeroAndHasRequestBody(t *testing.T) {
	m := dsl.PutMethod{}
	assert.False(t, m.IsNonZero())
	assert.False(t, m.HasRequestBody())

	m = dsl.PutMethod{RequestBody: "Request"}
	assert.True(t, m.IsNonZero())
	assert.True(t, m.HasRequestBody())

	m = dsl.PutMethod{Responses: dsl.Responses{200: "Response"}}
	assert.True(t, m.IsNonZero())
	assert.False(t, m.HasRequestBody())

	m = dsl.PutMethod{RequestBody: "Request", Responses: dsl.Responses{200: "Response"}}
	assert.True(t, m.IsNonZero())
	assert.True(t, m.HasRequestBody())
}

func TestDeleteMethodIsNonZero(t *testing.T) {
	m := dsl.DeleteMethod{}
	assert.False(t, m.IsNonZero())

	m = dsl.DeleteMethod{Responses: dsl.Responses{204: "Response"}}
	assert.True(t, m.IsNonZero())
}
