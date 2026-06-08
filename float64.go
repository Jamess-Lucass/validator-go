package validator

import (
	"fmt"
	"math"
	"reflect"
)

type Float64Rule struct {
	fieldBase[float64]
}

func (v *Validator) Float64(field *float64) *Float64Rule {
	r := &Float64Rule{fieldBase[float64]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func Float64() *Float64Rule {
	return &Float64Rule{}
}

func (r *Float64Rule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a number", coerceFloat64)
}

// coerceFloat64 converts any numeric value to float64.
func coerceFloat64(value any) (float64, bool) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), true
	default:
		return 0, false
	}
}

// Presence rules

func (r *Float64Rule) NotNil() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *Float64Rule) NotEmpty() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val != 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Comparison rules

func (r *Float64Rule) Lt(n float64) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val < n },
		message:  fmt.Sprintf("must be less than %v", n),
	})
	return r
}

func (r *Float64Rule) Lte(n float64) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val <= n },
		message:  fmt.Sprintf("must be less than or equal to %v", n),
	})
	return r
}

func (r *Float64Rule) Gt(n float64) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val > n },
		message:  fmt.Sprintf("must be greater than %v", n),
	})
	return r
}

func (r *Float64Rule) Gte(n float64) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val >= n },
		message:  fmt.Sprintf("must be greater than or equal to %v", n),
	})
	return r
}

// Sign rules

func (r *Float64Rule) Positive() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val > 0 },
		message:  msgPositive,
	})
	return r
}

func (r *Float64Rule) Negative() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val < 0 },
		message:  msgNegative,
	})
	return r
}

func (r *Float64Rule) Nonnegative() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val >= 0 },
		message:  msgNonnegative,
	})
	return r
}

func (r *Float64Rule) Nonpositive() *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool { return val <= 0 },
		message:  msgNonpositive,
	})
	return r
}

// Arithmetic rules

func (r *Float64Rule) MultipleOf(n float64) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{
		validate: func(val float64) bool {
			if n == 0 {
				return false
			}
			// Magnitude-relative tolerance so the check holds for large values.
			ratio := val / n
			rounded := math.Round(ratio)
			return math.Abs(ratio-rounded) <= 1e-9*math.Max(1, math.Abs(ratio))
		},
		message: fmt.Sprintf("must be a multiple of %v", n),
	})
	return r
}

// Custom rules

func (r *Float64Rule) Must(fn func(float64) bool) *Float64Rule {
	r.rules = append(r.rules, rule[float64]{validate: fn, message: msgNotValid})
	return r
}

func (r *Float64Rule) WithMessage(msg string) *Float64Rule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *Float64Rule) WithName(name string) *Float64Rule {
	r.name = name
	return r
}
