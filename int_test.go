package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestInt_Type(t *testing.T) {
	s := schema.Int()

	assert.True(t, s.Parse(0).IsValid())
	assert.True(t, s.Parse(int(uint(10))).IsValid())

	assert.False(t, s.Parse("123").IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0.015).IsValid())

	assert.False(t, s.Parse(uint(10)).IsValid())
	assert.False(t, s.Parse(uint8(10)).IsValid())
	assert.False(t, s.Parse(uint16(10)).IsValid())
	assert.False(t, s.Parse(uint32(10)).IsValid())
	assert.False(t, s.Parse(uint64(10)).IsValid())

	assert.False(t, s.Parse(int8(10)).IsValid())
	assert.False(t, s.Parse(int16(10)).IsValid())
	assert.False(t, s.Parse(int32(10)).IsValid())
	assert.False(t, s.Parse(int64(10)).IsValid())
}

func TestInt_Lt(t *testing.T) {
	s := schema.Int().Lt(5)

	assert.True(t, s.Parse(4).IsValid())
	assert.True(t, s.Parse(-5).IsValid())

	assert.False(t, s.Parse(5).IsValid())
	assert.False(t, s.Parse(500).IsValid())
}

func TestInt_Lte(t *testing.T) {
	s := schema.Int().Lte(5)

	assert.True(t, s.Parse(4).IsValid())
	assert.True(t, s.Parse(-5).IsValid())
	assert.True(t, s.Parse(5).IsValid())

	assert.False(t, s.Parse(6).IsValid())
	assert.False(t, s.Parse(500).IsValid())
}

func TestInt_Gt(t *testing.T) {
	s := schema.Int().Gt(5)

	assert.True(t, s.Parse(6).IsValid())
	assert.True(t, s.Parse(500).IsValid())

	assert.False(t, s.Parse(5).IsValid())
	assert.False(t, s.Parse(-5).IsValid())
}

func TestInt_Gte(t *testing.T) {
	s := schema.Int().Gte(5)

	assert.True(t, s.Parse(5).IsValid())
	assert.True(t, s.Parse(6).IsValid())
	assert.True(t, s.Parse(500).IsValid())

	assert.False(t, s.Parse(4).IsValid())
	assert.False(t, s.Parse(-5).IsValid())
}

func TestInt_Positive(t *testing.T) {
	s := schema.Int().Positive()

	assert.True(t, s.Parse(1).IsValid())
	assert.True(t, s.Parse(6).IsValid())
	assert.True(t, s.Parse(500).IsValid())

	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse(-5).IsValid())
}

func TestInt_Nonnegative(t *testing.T) {
	s := schema.Int().Nonnegative()

	assert.True(t, s.Parse(0).IsValid())
	assert.True(t, s.Parse(6).IsValid())
	assert.True(t, s.Parse(500).IsValid())

	assert.False(t, s.Parse(-5).IsValid())
}

func TestInt_Negative(t *testing.T) {
	s := schema.Int().Negative()

	assert.True(t, s.Parse(-1).IsValid())
	assert.True(t, s.Parse(-500).IsValid())

	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse(5).IsValid())
}

func TestInt_Nonpositive(t *testing.T) {
	s := schema.Int().Nonpositive()

	assert.True(t, s.Parse(0).IsValid())
	assert.True(t, s.Parse(-500).IsValid())

	assert.False(t, s.Parse(1).IsValid())
	assert.False(t, s.Parse(5).IsValid())
}

func TestInt_MultipleOf(t *testing.T) {
	s := schema.Int().MultipleOf(5)

	assert.True(t, s.Parse(0).IsValid())
	assert.True(t, s.Parse(-5).IsValid())
	assert.True(t, s.Parse(-25).IsValid())
	assert.True(t, s.Parse(50).IsValid())

	assert.False(t, s.Parse(1).IsValid())
	assert.False(t, s.Parse(3).IsValid())
}
