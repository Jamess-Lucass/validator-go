package validator

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// NumberRule validates a field of any real (integer or float) type, including
// named types like `type Age int`.
type NumberRule[T Real] struct {
	fieldBase[T]
}

// Number validates a numeric field of any real type, from int and float64 to
// named types like type Age int. The type is inferred from the field pointer:
//
//	validator.Number(v, &s.Age).Gte(0).Lte(150)
func Number[T Real](v *Validator, field *T) *NumberRule[T] {
	r := &NumberRule[T]{fieldBase[T]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use. Any numeric kind or numeric
// string is coerced into T; out-of-range values are rejected, not wrapped.
func (r *NumberRule[T]) Validate(value any) []ValidationError {
	typeErr := "must be a number"
	if isIntegerType[T]() {
		typeErr = "must be an integer"
	}
	return validateStandalone(r.rules, value, typeErr, coerceReal[T])
}

// Presence rules

// NotNil fails when the field pointer or standalone value is nil.
func (r *NumberRule[T]) NotNil() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on zero.
func (r *NumberRule[T]) NotEmpty() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val != 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Comparison rules

// Lt requires a value less than n.
func (r *NumberRule[T]) Lt(n T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val < n },
		message:  fmt.Sprintf("must be less than %v", n),
	})
	return r
}

// Lte requires a value less than or equal to n.
func (r *NumberRule[T]) Lte(n T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val <= n },
		message:  fmt.Sprintf("must be less than or equal to %v", n),
	})
	return r
}

// Gt requires a value greater than n.
func (r *NumberRule[T]) Gt(n T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val > n },
		message:  fmt.Sprintf("must be greater than %v", n),
	})
	return r
}

// Gte requires a value greater than or equal to n.
func (r *NumberRule[T]) Gte(n T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val >= n },
		message:  fmt.Sprintf("must be greater than or equal to %v", n),
	})
	return r
}

// Sign rules

// Positive requires a value greater than zero.
func (r *NumberRule[T]) Positive() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val > 0 },
		message:  msgPositive,
	})
	return r
}

// Negative requires a value less than zero.
func (r *NumberRule[T]) Negative() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val < 0 },
		message:  msgNegative,
	})
	return r
}

// Nonnegative requires a value of zero or more.
func (r *NumberRule[T]) Nonnegative() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val >= 0 },
		message:  msgNonnegative,
	})
	return r
}

// Nonpositive requires a value of zero or less.
func (r *NumberRule[T]) Nonpositive() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val <= 0 },
		message:  msgNonpositive,
	})
	return r
}

// OneOf requires the value to be one of the given values, e.g. the members of
// an integer enum.
func (r *NumberRule[T]) OneOf(values ...T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			for _, v := range values {
				if val == v {
					return true
				}
			}
			return false
		},
		message: msgOneOf(values, false),
	})
	return r
}

// Arithmetic rules

// MultipleOf requires the value to be a multiple of n. Exact for integers; a
// magnitude-relative tolerance for floats.
func (r *NumberRule[T]) MultipleOf(n T) *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return isMultipleOf(val, n) },
		message:  fmt.Sprintf("must be a multiple of %v", n),
	})
	return r
}

// Finite fails on NaN and ±Inf; integers always pass. The comparison rules
// don't reject NaN on their own, so chain this where that matters.
func (r *NumberRule[T]) Finite() *NumberRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			f := float64(val)
			return !math.IsNaN(f) && !math.IsInf(f, 0)
		},
		message: "must be a finite number",
	})
	return r
}

// Custom rules

// Must runs a custom check against the value. It panics if fn is nil.
func (r *NumberRule[T]) Must(fn func(T) bool) *NumberRule[T] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[T]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *NumberRule[T]) WithMessage(msg string) *NumberRule[T] {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *NumberRule[T]) WithName(name string) *NumberRule[T] {
	r.name = name
	return r
}

func isMultipleOf[T Real](val, n T) bool {
	vv := reflect.ValueOf(val)
	nv := reflect.ValueOf(n)

	switch vv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d := nv.Int()
		return d != 0 && vv.Int()%d == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		d := nv.Uint()
		return d != 0 && vv.Uint()%d == 0
	case reflect.Float32, reflect.Float64:
		d := nv.Float()
		if d == 0 {
			return false
		}
		ratio := vv.Float() / d
		return math.Abs(ratio-math.Round(ratio)) <= 1e-9*math.Max(1, math.Abs(ratio))
	default:
		return false
	}
}

// coerceReal converts any numeric kind, or a numeric string, into T. Values T
// can't hold exactly are rejected rather than wrapped or rounded: out of
// range, fractional floats into an integer type, integers above 2^53 into
// float64.
func coerceReal[T Real](value any) (T, bool) {
	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intToReal[T](rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintToReal[T](rv.Uint())
	case reflect.Float32, reflect.Float64:
		return floatToReal[T](rv.Float())
	case reflect.String:
		return parseReal[T](rv.String())
	default:
		var zero T
		return zero, false
	}
}

func intToReal[T Real](i int64) (T, bool) {
	var zero T
	out := reflect.New(reflect.TypeOf(zero)).Elem()

	switch out.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if out.OverflowInt(i) {
			return zero, false
		}
		out.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if i < 0 || out.OverflowUint(uint64(i)) {
			return zero, false
		}
		out.SetUint(uint64(i))
	default:
		f, ok := exactFloat(float64(i), out.Kind())
		if !ok || f >= math.MaxInt64 || int64(f) != i {
			return zero, false
		}
		out.SetFloat(f)
	}

	return out.Interface().(T), true
}

func uintToReal[T Real](u uint64) (T, bool) {
	var zero T
	out := reflect.New(reflect.TypeOf(zero)).Elem()

	switch out.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if u > math.MaxInt64 || out.OverflowInt(int64(u)) {
			return zero, false
		}
		out.SetInt(int64(u))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if out.OverflowUint(u) {
			return zero, false
		}
		out.SetUint(u)
	default:
		f, ok := exactFloat(float64(u), out.Kind())
		if !ok || f >= math.MaxUint64 || uint64(f) != u {
			return zero, false
		}
		out.SetFloat(f)
	}

	return out.Interface().(T), true
}

// exactFloat rejects a float that lost integer precision on the way to kind:
// float64 rounds above 2^53, float32 above 2^24.
func exactFloat(f float64, kind reflect.Kind) (float64, bool) {
	if kind == reflect.Float32 && float64(float32(f)) != f {
		return 0, false
	}
	return f, true
}

func floatToReal[T Real](f float64) (T, bool) {
	var zero T
	out := reflect.New(reflect.TypeOf(zero)).Elem()

	switch out.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Out-of-range float to int conversion is implementation-defined, so
		// check bounds in float form. MaxInt64 rounds up to 2^63, hence >=.
		if !isIntegralFloat(f) || f >= math.MaxInt64 || f < math.MinInt64 {
			return zero, false
		}
		i := int64(f)
		if out.OverflowInt(i) {
			return zero, false
		}
		out.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if !isIntegralFloat(f) || f < 0 || f >= math.MaxUint64 {
			return zero, false
		}
		u := uint64(f)
		if out.OverflowUint(u) {
			return zero, false
		}
		out.SetUint(u)
	default:
		// Matches parseReal, which rejects "NaN" and "Inf" strings.
		if math.IsNaN(f) || math.IsInf(f, 0) || out.OverflowFloat(f) {
			return zero, false
		}
		out.SetFloat(f)
	}

	return out.Interface().(T), true
}

func isIntegralFloat(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0) && f == math.Trunc(f)
}

func parseReal[T Real](s string) (T, bool) {
	var zero T

	if isIntegerType[T]() {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return intToReal[T](i)
		}
		// Values above MaxInt64 still fit unsigned targets.
		if u, err := strconv.ParseUint(s, 10, 64); err == nil {
			return uintToReal[T](u)
		}
		return zero, false
	}

	// Integer strings take the exact path; ParseFloat would round digits
	// above 2^53.
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return intToReal[T](i)
	}
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return uintToReal[T](u)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return zero, false
	}
	return floatToReal[T](f)
}

func isIntegerType[T Real]() bool {
	var zero T
	switch reflect.TypeOf(zero).Kind() {
	case reflect.Float32, reflect.Float64:
		return false
	default:
		return true
	}
}
