package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestArray_String_Type(t *testing.T) {
	s := schema.Array(schema.String())

	assert.True(t, s.Parse([]string{"one", "two"}).IsValid())
	assert.True(t, s.Parse([]string{""}).IsValid())
	assert.True(t, s.Parse([]string{}).IsValid())

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

func TestArray_Path(t *testing.T) {
	s := schema.Array(schema.String().Min(4)).Parse([]string{"one"})

	assert.Len(t, s.Errors, 1)
	assert.Equal(t, "0", s.Errors[0].Path)
}

func TestArray_String(t *testing.T) {
	s := schema.Array(schema.String().Min(4).StartsWith("a"))

	assert.True(t, s.Parse([]string{"aour", "aive"}).IsValid())
	assert.True(t, s.Parse([]string{"aour"}).IsValid())
	assert.True(t, s.Parse([]string{}).IsValid())

	assert.False(t, s.Parse([]string{"one"}).IsValid())
	assert.False(t, s.Parse([]string{"one", "four"}).IsValid())

	assert.Len(t, s.Parse([]string{"ane", "aour"}).Errors, 1)
	assert.Len(t, s.Parse([]string{"ane", "four"}).Errors, 2)
	assert.Len(t, s.Parse([]string{"one", "four"}).Errors, 3)
}

func TestArray_Min(t *testing.T) {
	s := schema.Array(schema.String()).Min(1)

	assert.True(t, s.Parse([]string{"aour", "aive"}).IsValid())
	assert.True(t, s.Parse([]string{"aour"}).IsValid())

	assert.False(t, s.Parse([]string{}).IsValid())
}

func TestArray_Max(t *testing.T) {
	s := schema.Array(schema.String()).Max(1)

	assert.True(t, s.Parse([]string{"aour"}).IsValid())
	assert.True(t, s.Parse([]string{}).IsValid())

	assert.False(t, s.Parse([]string{"aour", "aive"}).IsValid())
}

func TestArray_Object(t *testing.T) {
	type User struct {
		Firstname string
	}

	s := schema.Array(schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(4),
	}))

	val := []User{
		{Firstname: "1234"},
		{Firstname: "1234"},
	}

	assert.True(t, s.Parse(val).IsValid())
}

func TestArray_Object_Path(t *testing.T) {
	type User struct {
		Firstname string
	}

	s := schema.Array(schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(4),
	})).Parse([]User{
		{Firstname: "123"},
		{Firstname: "123"},
	})

	assert.Len(t, s.Errors, 2)
	assert.Equal(t, "0.Firstname", s.Errors[0].Path)
}
