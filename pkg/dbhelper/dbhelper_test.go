package dbhelper_test

import (
	"testing"

	"github.com/contracttests/broker/server/pkg/dbhelper"
	"github.com/stretchr/testify/assert"
)

func TestNullableStringReturnsNilForEmpty(t *testing.T) {
	assert.Nil(t, dbhelper.NullableString(""))
}

func TestNullableStringReturnsPointerForNonEmpty(t *testing.T) {
	got := dbhelper.NullableString("payments")

	if assert.NotNil(t, got) {
		assert.Equal(t, "payments", *got)
	}
}
