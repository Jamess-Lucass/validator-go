package validator

import (
	"math"
	"reflect"
	"strconv"
)

// ComplexRule validates a complex-number field (complex64 or complex128).
// Complex values aren't ordered, so only presence and custom rules apply.
type ComplexRule[T ComplexNumber] struct {
	fieldBase[T]
}

// Complex validates a complex-number field.
func Complex[T ComplexNumber](v *Validator, field *T) *ComplexRule[T] {
	r := &ComplexRule[T]{fieldBase[T]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use; complex values, real numbers,
// and parseable strings are accepted.
func (r *ComplexRule[T]) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a complex number", coerceComplex[T])
}

// NotNil fails when the field pointer or standalone value is nil.
func (r *ComplexRule[T]) NotNil() *ComplexRule[T] {
	r.rules = append(r.rules, rule[T]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on zero.
func (r *ComplexRule[T]) NotEmpty() *ComplexRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val != 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Must runs a custom check against the value. It panics if fn is nil.
func (r *ComplexRule[T]) Must(fn func(T) bool) *ComplexRule[T] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[T]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *ComplexRule[T]) WithMessage(msg string) *ComplexRule[T] {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *ComplexRule[T]) WithName(name string) *ComplexRule[T] {
	r.name = name
	return r
}

// coerceComplex converts a complex value, a real number, or a parseable
// string into T. NaN, ±Inf, and magnitudes that overflow complex64 are
// rejected, mirroring coerceReal.
func coerceComplex[T ComplexNumber](value any) (T, bool) {
	var zero T
	rv := reflect.ValueOf(value)

	var c complex128
	switch rv.Kind() {
	case reflect.Complex64, reflect.Complex128:
		c = rv.Complex()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		c = complex(float64(rv.Int()), 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		c = complex(float64(rv.Uint()), 0)
	case reflect.Float32, reflect.Float64:
		c = complex(rv.Float(), 0)
	case reflect.String:
		parsed, err := strconv.ParseComplex(rv.String(), 128)
		if err != nil {
			return zero, false
		}
		c = parsed
	default:
		return zero, false
	}

	out := T(c)
	if !isFiniteComplex(complex128(out)) {
		return zero, false
	}
	return out, true
}

// Narrowing a large complex128 to complex64 produces ±Inf, which this rejects.
func isFiniteComplex(c complex128) bool {
	re, im := real(c), imag(c)
	return !math.IsNaN(re) && !math.IsInf(re, 0) && !math.IsNaN(im) && !math.IsInf(im, 0)
}
