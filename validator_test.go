package validator_test

import (
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
	v.String(&s.Address.Street).NotEmpty()
	v.String(&s.Address.City).Min(3)
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
	v.String(&s.Name).NotEmpty()
	if s.Billing != nil {
		v.String(&s.Billing.Street).NotEmpty()
	}
	result := v.Validate()
	assert.True(t, result.IsValid())

	// When billing is present, validate it
	s.Billing = &nestedAddress{Street: "", City: ""}
	v, _ = validator.New(&s)
	if s.Billing != nil {
		v.String(&s.Billing.Street).NotEmpty()
		v.String(&s.Billing.City).NotEmpty()
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
	v.String(&s.FirstName).NotEmpty()
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
	v.String(s.Notes).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "unknown", result.Errors[0].Field)

	// WithName gives the nil field a stable name.
	v, _ = validator.New(&s)
	v.String(s.Notes).NotNil().WithName("notes")
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
