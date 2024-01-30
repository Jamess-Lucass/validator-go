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

func TestInt_Gt(t *testing.T) {
	s := schema.Int().Gt(5)

	assert.True(t, s.Parse(6).IsValid())
	assert.True(t, s.Parse(500).IsValid())

	assert.False(t, s.Parse(5).IsValid())
	assert.False(t, s.Parse(-5).IsValid())
}
