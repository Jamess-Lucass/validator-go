package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestBool_Type(t *testing.T) {
	s := schema.Bool()

	assert.True(t, s.Parse(true).IsValid())
	assert.True(t, s.Parse(false).IsValid())

	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse(1).IsValid())
	assert.False(t, s.Parse("true").IsValid())
	assert.False(t, s.Parse("false").IsValid())
	assert.False(t, s.Parse(uint(0)).IsValid())
}
