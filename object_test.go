package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

func TestNewObject_NilMap(t *testing.T) {
	// A nil map is usable: every key reads as missing, so rules skip unless
	// NotNil is chained. This is why NewObject has no error return.
	mv := validator.NewObject(nil)
	mv.Field("email", rule.String().Email())
	assert.True(t, mv.Validate().IsValid())

	mv = validator.NewObject(nil)
	mv.Field("email", rule.String().NotNil())
	res := mv.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "email", res.Errors[0].Field)
}

func TestMap_Basic(t *testing.T) {
	m := map[string]any{
		"email": "test@example.com",
		"age":   float64(20),
	}

	v := validator.NewObject(m)
	v.Field("email", rule.String().Email())
	v.Field("age", rule.Int().Gte(18))
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestMap_Failures(t *testing.T) {
	m := map[string]any{
		"email": "not-an-email",
		"age":   float64(10),
	}

	v := validator.NewObject(m)
	v.Field("email", rule.String().Email())
	v.Field("age", rule.Int().Gte(18))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 2)
}

func TestMap_TypeMismatch(t *testing.T) {
	m := map[string]any{
		"age": "not a number",
	}

	v := validator.NewObject(m)
	v.Field("age", rule.Int().Gte(18))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "age", result.Errors[0].Field)
	assert.Equal(t, "must be an integer", result.Errors[0].Message)
}

func TestMap_MissingKey(t *testing.T) {
	m := map[string]any{}

	// Missing key = nil, should skip
	v := validator.NewObject(m)
	v.Field("email", rule.String().Email())
	result := v.Validate()
	assert.True(t, result.IsValid())

	// Missing key with NotNil
	v = validator.NewObject(m)
	v.Field("email", rule.String().NotNil())
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestMap_Nested(t *testing.T) {
	m := map[string]any{
		"address": map[string]any{
			"city": "NY",
		},
	}

	v := validator.NewObject(m)
	v.Field("address", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("city", rule.String().Min(3))
	}))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "address.city", result.Errors[0].Field)
}

func TestMap_Slice(t *testing.T) {
	m := map[string]any{
		"tags": []any{"go", "a", "rust"},
	}

	v := validator.NewObject(m)
	v.Field("tags", validator.SliceOf(rule.String().Min(2)).Min(1))
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

	v := validator.NewObject(m)
	v.Field("users", validator.SliceOf(validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("name", rule.String().Min(2))
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

	v := validator.NewObject(m)
	v.Field("role", rule.String().NotEmpty())
	if m["role"] == "admin" {
		v.Field("admin_code", rule.String().NotEmpty())
	}
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "admin_code", result.Errors[0].Field)
}

func TestMap_ValidatesCurrentState(t *testing.T) {
	m := map[string]any{
		"email": "not-an-email",
	}

	v := validator.NewObject(m)
	v.Field("email", rule.String().Email())

	// Mutate the map after building rules but before validating
	m["email"] = "valid@example.com"

	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestObject_NilSkipped(t *testing.T) {
	v := validator.NewObject(map[string]any{})
	v.Field("address", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("city", rule.String().NotEmpty())
	}))
	assert.True(t, v.Validate().IsValid())
}

func TestObject_TypedNestedMap(t *testing.T) {
	payload := map[string]any{"address": map[string]string{"city": ""}}
	v := validator.NewObject(payload)
	v.Field("address", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("city", rule.String().NotEmpty())
	}))
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "address.city", res.Errors[0].Field)
}

func TestMap_NestedNotNil(t *testing.T) {
	m := map[string]any{}

	v := validator.NewObject(m)
	v.Field("profile", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("bio", rule.String().NotEmpty())
	}).NotNil())
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "profile", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestObject_NilCallbackPanics(t *testing.T) {
	// A nil callback is a programmer error caught at construction, not a nil
	// deref later when a value happens to reach Validate.
	assert.Panics(t, func() { validator.Object(nil) })
}

func TestField_NilRulePanics(t *testing.T) {
	mv := validator.NewObject(map[string]any{"name": "x"})
	assert.Panics(t, func() { mv.Field("name", nil) })

	// A typed nil boxed into the Rule interface is caught the same way.
	var typedNil *validator.StringRule[string]
	assert.Panics(t, func() { mv.Field("name", typedNil) })
}
