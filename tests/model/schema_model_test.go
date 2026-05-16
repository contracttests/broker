package model_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewSchema(t *testing.T) {
	s := model.NewSchema()
	assert.NotNil(t, s.Properties)
	assert.Empty(t, s.Properties)
}

func TestSchemaAddProperty(t *testing.T) {
	s := model.NewSchema()
	property := model.NewProperty("root", "object", false)
	s.AddProperty(property)

	assert.Equal(t, model.SchemaProperties{"root": property}, s.Properties)
}
