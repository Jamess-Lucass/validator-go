package validator

import "reflect"

// BoolRule validates a bool field or value. Named bool types work too.
type BoolRule[T ~bool] struct {
	fieldBase[T]
}

// Bool validates a bool field, plain or named.
func Bool[T ~bool](v *Validator, field *T) *BoolRule[T] {
	r := &BoolRule[T]{fieldBase[T]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use; the value must be a bool or a
// named bool type.
func (r *BoolRule[T]) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a boolean", func(v any) (T, bool) {
		if b, ok := v.(T); ok {
			return b, true
		}
		if rv := reflect.ValueOf(v); rv.Kind() == reflect.Bool {
			return T(rv.Bool()), true
		}
		return false, false
	})
}

// NotNil fails when the field pointer or standalone value is nil.
func (r *BoolRule[T]) NotNil() *BoolRule[T] {
	r.rules = append(r.rules, rule[T]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on false.
func (r *BoolRule[T]) NotEmpty() *BoolRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return bool(val) },
		message:  msgNotEmpty,
	})
	return r
}

// Must runs a custom check against the value. It panics if fn is nil.
func (r *BoolRule[T]) Must(fn func(T) bool) *BoolRule[T] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[T]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *BoolRule[T]) WithMessage(msg string) *BoolRule[T] {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *BoolRule[T]) WithName(name string) *BoolRule[T] {
	r.name = name
	return r
}
