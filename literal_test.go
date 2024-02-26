package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestLiteral_String(t *testing.T) {
	s := schema.Literal("test")

	assert.True(t, s.Parse("test").IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
}

func TestLiteral_Int(t *testing.T) {
	s := schema.Literal(10)

	assert.True(t, s.Parse(10).IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
}

func TestLiteral_Struct(t *testing.T) {
	type User struct {
		FirstName string
	}

	s := schema.Literal(User{FirstName: "John"})

	assert.True(t, s.Parse(User{FirstName: "John"}).IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse(User{FirstName: "John1"}).IsValid())
	assert.False(t, s.Parse(User{}).IsValid())
}
