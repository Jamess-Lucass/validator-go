package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestFloat64_Type(t *testing.T) {
	s := schema.Float64()

	assert.True(t, s.Parse(float64(0)).IsValid())
	assert.True(t, s.Parse(0.0).IsValid())
	assert.True(t, s.Parse(0.01).IsValid())
	assert.True(t, s.Parse(0.015).IsValid())
	assert.True(t, s.Parse(float64(uint(10))).IsValid())

	assert.False(t, s.Parse("123").IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]float64{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]float64{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse(10).IsValid())

	assert.False(t, s.Parse(int(10)).IsValid())
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

func TestFloat64_Lt(t *testing.T) {
	s := schema.Float64().Lt(5)

	assert.True(t, s.Parse(2.5).IsValid())
	assert.True(t, s.Parse(float64(4)).IsValid())
	assert.True(t, s.Parse(float64(-5)).IsValid())

	assert.False(t, s.Parse(float64(5)).IsValid())
	assert.False(t, s.Parse(float64(500)).IsValid())
}

func TestFloat64_Lte(t *testing.T) {
	s := schema.Float64().Lte(5)

	assert.True(t, s.Parse(5.00).IsValid())
	assert.True(t, s.Parse(-7.45).IsValid())
	assert.True(t, s.Parse(4.5).IsValid())
	assert.True(t, s.Parse(float64(4)).IsValid())
	assert.True(t, s.Parse(float64(-5)).IsValid())
	assert.True(t, s.Parse(float64(5)).IsValid())

	assert.False(t, s.Parse(float64(6)).IsValid())
	assert.False(t, s.Parse(float64(500)).IsValid())
}

func TestFloat64_Gt(t *testing.T) {
	s := schema.Float64().Gt(5)

	assert.True(t, s.Parse(6.5).IsValid())
	assert.True(t, s.Parse(float64(6)).IsValid())
	assert.True(t, s.Parse(float64(500)).IsValid())

	assert.False(t, s.Parse(float64(5)).IsValid())
	assert.False(t, s.Parse(float64(-5)).IsValid())
}

func TestFloat64_Gte(t *testing.T) {
	s := schema.Float64().Gte(5)

	assert.True(t, s.Parse(5.00).IsValid())
	assert.True(t, s.Parse(5.01).IsValid())
	assert.True(t, s.Parse(float64(5)).IsValid())
	assert.True(t, s.Parse(float64(6)).IsValid())
	assert.True(t, s.Parse(float64(500)).IsValid())

	assert.False(t, s.Parse(4.99).IsValid())
	assert.False(t, s.Parse(float64(4)).IsValid())
	assert.False(t, s.Parse(float64(-5)).IsValid())
}

func TestFloat64_Positive(t *testing.T) {
	s := schema.Float64().Positive()

	assert.True(t, s.Parse(0.01).IsValid())
	assert.True(t, s.Parse(float64(1)).IsValid())
	assert.True(t, s.Parse(float64(6)).IsValid())
	assert.True(t, s.Parse(float64(500)).IsValid())

	assert.False(t, s.Parse(-0.01).IsValid())
	assert.False(t, s.Parse(float64(0)).IsValid())
	assert.False(t, s.Parse(float64(-5)).IsValid())
}

func TestFloat64_Nonnegative(t *testing.T) {
	s := schema.Float64().Nonnegative()

	assert.True(t, s.Parse(0.01).IsValid())
	assert.True(t, s.Parse(0.00).IsValid())
	assert.True(t, s.Parse(float64(0)).IsValid())
	assert.True(t, s.Parse(float64(6)).IsValid())
	assert.True(t, s.Parse(float64(500)).IsValid())

	assert.False(t, s.Parse(float64(-5)).IsValid())
	assert.False(t, s.Parse(-0.01).IsValid())
}

func TestFloat64_Negative(t *testing.T) {
	s := schema.Float64().Negative()

	assert.True(t, s.Parse(-0.01).IsValid())
	assert.True(t, s.Parse(float64(-1)).IsValid())
	assert.True(t, s.Parse(float64(-500)).IsValid())

	assert.False(t, s.Parse(0.00).IsValid())
	assert.False(t, s.Parse(0.01).IsValid())
	assert.False(t, s.Parse(-0.00).IsValid())
	assert.False(t, s.Parse(float64(0)).IsValid())
	assert.False(t, s.Parse(float64(5)).IsValid())
}

func TestFloat64_Nonpositive(t *testing.T) {
	s := schema.Float64().Nonpositive()

	assert.True(t, s.Parse(0.00).IsValid())
	assert.True(t, s.Parse(-0.01).IsValid())
	assert.True(t, s.Parse(float64(0)).IsValid())
	assert.True(t, s.Parse(float64(-500)).IsValid())

	assert.False(t, s.Parse(0.01).IsValid())
	assert.False(t, s.Parse(float64(1)).IsValid())
	assert.False(t, s.Parse(float64(5)).IsValid())
}

func TestFloat64_MultipleOf(t *testing.T) {
	s := schema.Float64().MultipleOf(5)

	assert.True(t, s.Parse(0.00).IsValid())
	assert.True(t, s.Parse(10.00).IsValid())
	assert.True(t, s.Parse(float64(0)).IsValid())
	assert.True(t, s.Parse(float64(-5)).IsValid())
	assert.True(t, s.Parse(float64(-25)).IsValid())
	assert.True(t, s.Parse(float64(50)).IsValid())

	assert.False(t, s.Parse(10.01).IsValid())
	assert.False(t, s.Parse(float64(1)).IsValid())
	assert.False(t, s.Parse(float64(3)).IsValid())
}
