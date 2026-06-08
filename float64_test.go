package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type float64TestStruct struct {
	Price float64  `json:"price"`
	Rate  *float64 `json:"rate"`
}

func TestFloat64_NotNil_WithName(t *testing.T) {
	s := float64TestStruct{Rate: nil}
	v, _ := validator.New(&s)
	v.Float64(s.Rate).NotNil().WithName("rate")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "rate", result.Errors[0].Field)
}

func TestFloat64_Gte(t *testing.T) {
	s := float64TestStruct{Price: -1.0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Gte(0)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Gte(0)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_Positive(t *testing.T) {
	s := float64TestStruct{Price: 0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Positive()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 0.01
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Positive()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_NotEmpty(t *testing.T) {
	s := float64TestStruct{Price: 0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 5.5
	v, _ = validator.New(&s)
	v.Float64(&s.Price).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_Lt(t *testing.T) {
	s := float64TestStruct{Price: 5.0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Lt(5.0)
	assert.False(t, v.Validate().IsValid())

	s.Price = 4.9
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Lt(5.0)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Gt(t *testing.T) {
	s := float64TestStruct{Price: 5.0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Gt(5.0)
	assert.False(t, v.Validate().IsValid())

	s.Price = 5.1
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Gt(5.0)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Negative(t *testing.T) {
	s := float64TestStruct{Price: 1.0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Negative()
	assert.False(t, v.Validate().IsValid())

	s.Price = -1.0
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Negative()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Nonnegative(t *testing.T) {
	s := float64TestStruct{Price: -0.1}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Nonnegative()
	assert.False(t, v.Validate().IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Nonnegative()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Nonpositive(t *testing.T) {
	s := float64TestStruct{Price: 0.1}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Nonpositive()
	assert.False(t, v.Validate().IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	v.Float64(&s.Price).Nonpositive()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_MultipleOf(t *testing.T) {
	s := float64TestStruct{Price: 0.7}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).MultipleOf(0.25)
	assert.False(t, v.Validate().IsValid())

	s.Price = 0.75
	v, _ = validator.New(&s)
	v.Float64(&s.Price).MultipleOf(0.25)
	assert.True(t, v.Validate().IsValid())

	// Large magnitude stays robust (a fixed epsilon would mis-handle this).
	s.Price = 1e12
	v, _ = validator.New(&s)
	v.Float64(&s.Price).MultipleOf(0.25)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Must(t *testing.T) {
	s := float64TestStruct{Price: 3.0}
	v, _ := validator.New(&s)
	v.Float64(&s.Price).Must(func(val float64) bool {
		return val >= 1.0 && val <= 2.0
	}).WithMessage("must be between 1 and 2")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be between 1 and 2", result.Errors[0].Message)
}

func TestFloat64_Standalone(t *testing.T) {
	rule := validator.Float64().Gte(0).Lte(100)

	errs := rule.Validate(50.5)
	assert.Empty(t, errs)

	errs = rule.Validate(-1.0)
	assert.Len(t, errs, 1)

	errs = rule.Validate("string")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a number", errs[0].Message)
}

func TestFloat64_Coercion_AllNumericKinds(t *testing.T) {
	rule := validator.Float64().Gte(0)

	assert.Empty(t, rule.Validate(float32(1.5)))
	assert.Empty(t, rule.Validate(float64(1.5)))
	assert.Empty(t, rule.Validate(int(2)))
	assert.Empty(t, rule.Validate(int64(2)))
	assert.Empty(t, rule.Validate(uint(2)))

	errs := rule.Validate(true)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a number", errs[0].Message)
}
