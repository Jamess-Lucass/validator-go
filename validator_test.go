package validator_test

import (
	"net/url"
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type nestedAddress struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type nestedTestStruct struct {
	Name    string         `json:"name"`
	Address nestedAddress  `json:"address"`
	Billing *nestedAddress `json:"billing"`
}

func TestNested_DotNotation(t *testing.T) {
	s := nestedTestStruct{
		Name:    "John",
		Address: nestedAddress{Street: "", City: "NY"},
	}
	v, _ := validator.New(&s)
	validator.String(v, &s.Address.Street).NotEmpty()
	validator.String(v, &s.Address.City).Min(3)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "address.street", result.Errors[0].Field)
	assert.Equal(t, "address.city", result.Errors[1].Field)
}

func TestNested_NullableStruct(t *testing.T) {
	s := nestedTestStruct{
		Name:    "John",
		Address: nestedAddress{Street: "123 Main", City: "NYC"},
		Billing: nil,
	}

	// When billing is nil, we just don't validate it (plain if)
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).NotEmpty()
	if s.Billing != nil {
		validator.String(v, &s.Billing.Street).NotEmpty()
	}
	result := v.Validate()
	assert.True(t, result.IsValid())

	// When billing is present, validate it
	s.Billing = &nestedAddress{Street: "", City: ""}
	v, _ = validator.New(&s)
	if s.Billing != nil {
		validator.String(v, &s.Billing.Street).NotEmpty()
		validator.String(v, &s.Billing.City).NotEmpty()
	}
	result = v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "billing.street", result.Errors[0].Field)
	assert.Equal(t, "billing.city", result.Errors[1].Field)
}

func TestFieldName_FallbackToGoName(t *testing.T) {
	type noTags struct {
		FirstName string
	}
	s := noTags{FirstName: ""}
	v, _ := validator.New(&s)
	validator.String(v, &s.FirstName).NotEmpty()
	result := v.Validate()
	assert.Equal(t, "FirstName", result.Errors[0].Field)
}

func TestNew_InvalidArg(t *testing.T) {
	cases := []struct {
		name string
		arg  any
	}{
		{"non-pointer", nestedTestStruct{}},
		{"nil interface", nil},
		{"nil pointer", (*nestedTestStruct)(nil)},
		{"pointer to non-struct", new(int)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := validator.New(tc.arg)
			assert.Error(t, err)
			assert.Nil(t, v)
		})
	}
}

func TestNew_Valid(t *testing.T) {
	v, err := validator.New(&nestedTestStruct{})
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestWithName_NilNullablePointer(t *testing.T) {
	type withPtr struct {
		Notes *string `json:"notes"`
	}
	s := withPtr{Notes: nil}

	// Without WithName, a nil pointer cannot be resolved by reflection.
	v, _ := validator.New(&s)
	validator.String(v, s.Notes).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "unknown", result.Errors[0].Field)

	// WithName gives the nil field a stable name.
	v, _ = validator.New(&s)
	validator.String(v, s.Notes).NotNil().WithName("notes")
	result = v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "notes", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestValidationError_Error(t *testing.T) {
	withField := validator.ValidationError{Field: "name", Message: "must not be empty"}
	assert.Equal(t, "name: must not be empty", withField.Error())

	noField := validator.ValidationError{Message: "must not be empty"}
	assert.Equal(t, "must not be empty", noField.Error())
}

func TestNested_InlineStructField(t *testing.T) {
	type S struct {
		Meta struct {
			Note string `json:"note"`
		} `json:"meta"`
	}
	var s S
	v, _ := validator.New(&s)
	validator.String(v, &s.Meta.Note).NotEmpty()
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "meta.note", res.Errors[0].Field)
}

func TestNested_CrossPackageStruct(t *testing.T) {
	// A nested struct from another package resolves as long as it has
	// exported fields. url.URL stands in for a typical shared model type.
	type endpoint struct {
		URL url.URL `json:"url"`
	}
	e := endpoint{}
	v, _ := validator.New(&e)
	validator.String(v, &e.URL.Host).NotEmpty()
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "url.Host", res.Errors[0].Field)
}

func TestResolve_CyclicStructNoCrash(t *testing.T) {
	type Node struct {
		Parent *Node  `json:"-"`
		Name   string `json:"name"`
	}
	a := &Node{Name: ""}
	b := &Node{Name: "x", Parent: a}
	a.Parent = b // a <-> b cycle

	v, _ := validator.New(a)
	validator.String(v, &a.Name).NotEmpty()
	res := v.Validate()

	// Resolves the direct field without recursing into the cycle (no stack overflow).
	assert.False(t, res.IsValid())
	assert.Equal(t, "name", res.Errors[0].Field)
}
