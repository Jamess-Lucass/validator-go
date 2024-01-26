package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestInt_Type(t *testing.T) {
	s := schema.Int()

	assert.True(t, s.Parse(0))
	assert.True(t, s.Parse(int(uint(10))))

	assert.False(t, s.Parse("123"))
	assert.False(t, s.Parse(nil))
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}))
	assert.False(t, s.Parse([]int{1, 2, 3}))
	assert.False(t, s.Parse(0.015))

	assert.False(t, s.Parse(uint(10)))
	assert.False(t, s.Parse(uint8(10)))
	assert.False(t, s.Parse(uint16(10)))
	assert.False(t, s.Parse(uint32(10)))
	assert.False(t, s.Parse(uint64(10)))

	assert.False(t, s.Parse(int8(10)))
	assert.False(t, s.Parse(int16(10)))
	assert.False(t, s.Parse(int32(10)))
	assert.False(t, s.Parse(int64(10)))
}

func TestInt_LessThan(t *testing.T) {
	s := schema.Int().LessThan(5)

	assert.True(t, s.Parse(4))
	assert.True(t, s.Parse(-5))

	assert.False(t, s.Parse(5))
	assert.False(t, s.Parse(500))
}
