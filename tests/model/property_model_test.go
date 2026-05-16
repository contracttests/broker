package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewProperty(t *testing.T) {
	p := model.NewProperty("root.name", "string", true)
	assert.Equal(t, model.Property{Path: "root.name", Type: "string", Optional: true}, p)
}

func TestPropertyIsSame(t *testing.T) {
	a := model.NewProperty("root.id", "string", false)
	b := model.NewProperty("root.id", "string", false)
	assert.True(t, a.IsSame(&b))

	differentPath := model.NewProperty("root.name", "string", false)
	assert.False(t, a.IsSame(&differentPath))

	differentType := model.NewProperty("root.id", "integer", false)
	assert.False(t, a.IsSame(&differentType))

	differentOptional := model.NewProperty("root.id", "string", true)
	assert.False(t, a.IsSame(&differentOptional))
}
