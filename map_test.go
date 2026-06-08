package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestMap_Basic(t *testing.T) {
	m := map[string]any{
		"email": "test@example.com",
		"age":   float64(20),
	}

	v := validator.NewMap(m)
	v.Field("email", validator.String().Email())
	v.Field("age", validator.Int().Gte(18))
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestMap_Failures(t *testing.T) {
	m := map[string]any{
		"email": "not-an-email",
		"age":   float64(10),
	}

	v := validator.NewMap(m)
	v.Field("email", validator.String().Email())
	v.Field("age", validator.Int().Gte(18))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 2)
}

func TestMap_TypeMismatch(t *testing.T) {
	m := map[string]any{
		"age": "not a number",
	}

	v := validator.NewMap(m)
	v.Field("age", validator.Int().Gte(18))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "age", result.Errors[0].Field)
	assert.Equal(t, "must be an integer", result.Errors[0].Message)
}

func TestMap_MissingKey(t *testing.T) {
	m := map[string]any{}

	// Missing key = nil, should skip
	v := validator.NewMap(m)
	v.Field("email", validator.String().Email())
	result := v.Validate()
	assert.True(t, result.IsValid())

	// Missing key with NotNil
	v = validator.NewMap(m)
	v.Field("email", validator.String().NotNil())
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestMap_Nested(t *testing.T) {
	m := map[string]any{
		"address": map[string]any{
			"city": "NY",
		},
	}

	v := validator.NewMap(m)
	v.Field("address", validator.Map(func(mv *validator.MapV) {
		mv.Field("city", validator.String().Min(3))
	}))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "address.city", result.Errors[0].Field)
}

func TestMap_Slice(t *testing.T) {
	m := map[string]any{
		"tags": []any{"go", "a", "rust"},
	}

	v := validator.NewMap(m)
	v.Field("tags", validator.SliceOf(validator.String().Min(2)).Min(1))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "tags[1]", result.Errors[0].Field)
	assert.Equal(t, "must be at least 2 characters", result.Errors[0].Message)
}

func TestMap_SliceOfNestedMaps(t *testing.T) {
	m := map[string]any{
		"users": []any{
			map[string]any{"name": "ok"},
			map[string]any{"name": "x"},
		},
	}

	v := validator.NewMap(m)
	v.Field("users", validator.SliceOf(validator.Map(func(mv *validator.MapV) {
		mv.Field("name", validator.String().Min(2))
	})))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "users[1].name", result.Errors[0].Field)
	assert.Equal(t, "must be at least 2 characters", result.Errors[0].Message)
}

func TestMap_Conditional(t *testing.T) {
	m := map[string]any{
		"role":       "admin",
		"admin_code": "",
	}

	v := validator.NewMap(m)
	v.Field("role", validator.String().NotEmpty())
	if m["role"] == "admin" {
		v.Field("admin_code", validator.String().NotEmpty())
	}
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "admin_code", result.Errors[0].Field)
}

func TestMap_ValidatesCurrentState(t *testing.T) {
	m := map[string]any{
		"email": "not-an-email",
	}

	v := validator.NewMap(m)
	v.Field("email", validator.String().Email())

	// Mutate the map after building rules but before validating
	m["email"] = "valid@example.com"

	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestMap_NestedNotNil(t *testing.T) {
	m := map[string]any{}

	v := validator.NewMap(m)
	v.Field("profile", validator.Map(func(mv *validator.MapV) {
		mv.Field("bio", validator.String().NotEmpty())
	}).NotNil())
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "profile", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}
