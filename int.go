package validator

import (
	"fmt"
	"math"
	"reflect"
)

type IntRule struct {
	fieldBase[int]
}

func (v *Validator) Int(field *int) *IntRule {
	r := &IntRule{fieldBase[int]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func Int() *IntRule { return &IntRule{} }

func (r *IntRule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be an integer", coerceInt)
}

// coerceInt converts any integer-valued number to int. A float with a
// fractional part is rejected.
func coerceInt(value any) (int, bool) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := rv.Uint()
		if u > uint64(math.MaxInt) {
			return 0, false
		}
		return int(u), true
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) || f != math.Floor(f) {
			return 0, false
		}
		return int(f), true
	default:
		return 0, false
	}
}

// Presence rules

func (r *IntRule) NotNil() *IntRule {
	r.rules = append(r.rules, rule[int]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *IntRule) NotEmpty() *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val != 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Comparison rules

func (r *IntRule) Lt(n int) *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val < n },
		message:  fmt.Sprintf("must be less than %d", n),
	})
	return r
}

func (r *IntRule) Lte(n int) *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val <= n },
		message:  fmt.Sprintf("must be less than or equal to %d", n),
	})
	return r
}

func (r *IntRule) Gt(n int) *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val > n },
		message:  fmt.Sprintf("must be greater than %d", n),
	})
	return r
}

func (r *IntRule) Gte(n int) *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val >= n },
		message:  fmt.Sprintf("must be greater than or equal to %d", n),
	})
	return r
}

// Sign rules

func (r *IntRule) Positive() *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val > 0 },
		message:  msgPositive,
	})
	return r
}

func (r *IntRule) Negative() *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val < 0 },
		message:  msgNegative,
	})
	return r
}

func (r *IntRule) Nonnegative() *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val >= 0 },
		message:  msgNonnegative,
	})
	return r
}

func (r *IntRule) Nonpositive() *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool { return val <= 0 },
		message:  msgNonpositive,
	})
	return r
}

// Arithmetic rules

func (r *IntRule) MultipleOf(n int) *IntRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(val int) bool {
			if n == 0 {
				return false
			}
			return val%n == 0
		},
		message: fmt.Sprintf("must be a multiple of %d", n),
	})
	return r
}

// Custom rules

func (r *IntRule) Must(fn func(int) bool) *IntRule {
	r.rules = append(r.rules, rule[int]{validate: fn, message: msgNotValid})
	return r
}

func (r *IntRule) WithMessage(msg string) *IntRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *IntRule) WithName(name string) *IntRule {
	r.name = name
	return r
}
