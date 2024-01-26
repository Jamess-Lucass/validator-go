package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestString_Type(t *testing.T) {
	s := schema.String()

	assert.True(t, s.Parse("string"))

	assert.False(t, s.Parse(123))
	assert.False(t, s.Parse(nil))
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}))
	assert.False(t, s.Parse([]int{1, 2, 3}))
	assert.False(t, s.Parse(0))
}

func TestString_Min(t *testing.T) {
	s := schema.String().Min(5)

	assert.True(t, s.Parse("12345"))

	assert.False(t, s.Parse("1234"))
}
