package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestObject_Type(t *testing.T) {
	type User struct {
		Firstname string
		Lastname  string
	}

	s := schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(2),
		"Lastname":  schema.String().StartsWith("d"),
	})

	assert.True(t, s.Parse(User{Firstname: "john", Lastname: "doe"}).IsValid())

	assert.False(t, s.Parse(User{Firstname: "john", Lastname: ""}).IsValid())
	assert.False(t, s.Parse(User{Firstname: "", Lastname: "doe"}).IsValid())
	assert.False(t, s.Parse(User{Firstname: "", Lastname: ""}).IsValid())
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

func TestObject_Path(t *testing.T) {
	type User struct {
		Firstname string
		Lastname  string
	}

	s := schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(2),
		"Lastname":  schema.String().StartsWith("d"),
	}).Parse(User{Firstname: "", Lastname: ""})

	assert.Len(t, s.Errors, 2)
	assert.Equal(t, "Firstname", s.Errors[0].Path)
	assert.Equal(t, "Lastname", s.Errors[1].Path)
}

func TestObject_ArrayObject_Path(t *testing.T) {
	type Address struct {
		Postcode string
	}

	type User struct {
		Addresses []Address
	}

	s := schema.Object(map[string]schema.ISchema{
		"Addresses": schema.Array(schema.Object(map[string]schema.ISchema{
			"Postcode": schema.String().Min(4),
		})),
	}).Parse(User{
		Addresses: []Address{{Postcode: "123"}},
	})

	assert.Len(t, s.Errors, 1)
	assert.Equal(t, "Addresses.0.Postcode", s.Errors[0].Path)
}
