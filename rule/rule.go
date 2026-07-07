// Package rule provides standalone rule values: reusable validation chains
// that aren't bound to a struct field. They plug into NewObject fields,
// EachValue, EachKey, SliceOf, and MapOf, and coerce their input where that
// makes sense (numbers from any numeric kind or string, times from RFC3339
// strings). Error names come from the map key or element index, so WithName
// has no effect on a standalone rule.
//
//	mv := validator.NewObject(payload)
//	mv.Field("email", rule.String().NotEmpty().Email())
//	mv.Field("age", rule.Int().Gte(18))
package rule

import (
	validator "github.com/Jamess-Lucass/validator-go"
)

// String returns a string rule.
func String() *validator.StringRule[string] {
	return &validator.StringRule[string]{}
}

// Bool returns a bool rule.
func Bool() *validator.BoolRule[bool] {
	return &validator.BoolRule[bool]{}
}

// Time returns a time.Time rule. Strings are parsed as RFC3339.
func Time() *validator.TimeRule {
	return &validator.TimeRule{}
}

// UUID returns a uuid.UUID rule. Strings are parsed.
func UUID() *validator.UUIDRule {
	return &validator.UUIDRule{}
}

// Int returns an int rule. Other numeric kinds and numeric strings are
// coerced, rejecting fractions and out-of-range values.
func Int() *validator.NumberRule[int] {
	return &validator.NumberRule[int]{}
}

// Float64 returns a float64 rule.
func Float64() *validator.NumberRule[float64] {
	return &validator.NumberRule[float64]{}
}

// Number returns a rule for any real numeric type. The type parameter sets
// the coercion target, so rule.Number[uint8]() rejects 300 and -1.
func Number[T validator.Real]() *validator.NumberRule[T] {
	return &validator.NumberRule[T]{}
}

// Complex128 returns a complex128 rule. Real numbers and parseable strings
// are accepted.
func Complex128() *validator.ComplexRule[complex128] {
	return &validator.ComplexRule[complex128]{}
}

// Complex64 returns a complex64 rule.
func Complex64() *validator.ComplexRule[complex64] {
	return &validator.ComplexRule[complex64]{}
}
