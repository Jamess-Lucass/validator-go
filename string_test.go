package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestString_Type(t *testing.T) {
	s := schema.String()

	assert.True(t, s.Parse("string").IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
}

func TestString_Min(t *testing.T) {
	s := schema.String().Min(5)

	assert.True(t, s.Parse("12345").IsValid())
	assert.True(t, s.Parse("123456").IsValid())

	assert.False(t, s.Parse("1234").IsValid())
}

func TestString_Max(t *testing.T) {
	s := schema.String().Max(5)

	assert.True(t, s.Parse("12345").IsValid())
	assert.True(t, s.Parse("1234").IsValid())

	assert.False(t, s.Parse("123456").IsValid())
}
