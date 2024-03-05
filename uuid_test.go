package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID_Type(t *testing.T) {
	s := schema.UUID()

	assert.True(t, s.Parse(uuid.New()).IsValid())
	assert.True(t, s.Parse(uuid.MustParse("00000000-0000-0000-0000-000000000000")).IsValid())
	assert.True(t, s.Parse(uuid.MustParse("6e3c7cd3-fc85-4bd6-ab47-1fc0f236a774")).IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
	assert.False(t, s.Parse("57c6b6aa-211a-4b49-a012-3fd9b4a4ea2d").IsValid())
	assert.False(t, s.Parse("db9fb12c-daea-11ee-a506-0242ac120002").IsValid())
	assert.False(t, s.Parse("018e0e8f-b1d9-7503-ac4a-49b18a95be69").IsValid())
	assert.False(t, s.Parse("00000000-0000-0000-0000-000000000000").IsValid())
}
