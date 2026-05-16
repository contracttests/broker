package dsl_test

import (
	"testing"

	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestPropertyPathString(t *testing.T) {
	p := dsl.NewPropertyPath("root")
	assert.Equal(t, "root", p.String())
}

func TestPropertyPathAppend(t *testing.T) {
	p := dsl.NewPropertyPath("root")
	p = p.Append("name")
	assert.Equal(t, "root.name", p.String())
}

func TestPropertyPathAppendArray(t *testing.T) {
	p := dsl.NewPropertyPath("root.tags")
	p = p.AppendArray()
	assert.Equal(t, "root.tags[]", p.String())
}