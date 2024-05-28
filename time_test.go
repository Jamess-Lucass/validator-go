package schema_test

import (
	"testing"
	"time"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestTime_Type(t *testing.T) {
	s := schema.Time()

	assert.True(t, s.Parse(time.Now()).IsValid())
	assert.True(t, s.Parse(time.Now().UTC()).IsValid())
	assert.True(t, s.Parse(time.Now().UTC().Add(time.Hour*2)).IsValid())

	assert.False(t, s.Parse(time.Now().String()).IsValid())
	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
}
