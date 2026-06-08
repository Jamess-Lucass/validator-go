package validator

import (
	"fmt"
	"reflect"
)

type SliceRule[T any] struct {
	fieldBase[[]T]
	eachFn   func(item *T, sv *Validator)
	eachRule Rule
}

// Slice is a free function because Go methods cannot have their own type parameters.
func Slice[T any](v *Validator, field *[]T) *SliceRule[T] {
	r := &SliceRule[T]{fieldBase: fieldBase[[]T]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func (r *SliceRule[T]) execute(prefix string) []ValidationError {
	fieldName := r.resolveName(prefix)

	if r.fieldPtr == nil {
		return executeRules(r.rules, nil, fieldName, true)
	}

	val := *r.fieldPtr

	var errs []ValidationError
	// A nil slice value behind a non-nil pointer fails NotNil, but length and
	// element rules still run against it as an empty slice.
	if val == nil {
		errs = append(errs, notNilErrors(r.rules, fieldName)...)
	}
	errs = append(errs, executeRules(r.rules, val, fieldName, false)...)

	if r.eachFn != nil {
		for i := range val {
			item := &val[i]
			sv := &Validator{structPtr: item, prefix: fmt.Sprintf("%s[%d]", fieldName, i)}
			r.eachFn(item, sv)
			result := sv.Validate()
			errs = append(errs, result.Errors...)
		}
	}

	errs = append(errs, r.eachValueErrors(val, fieldName)...)

	return errs
}

// Validate implements Rule for standalone/map-based slice validation.
// Note: Each(fn) callbacks are only executed in struct context (via execute),
// not in standalone mode, because they require a struct-based *Validator for field
// name resolution.
func (r *SliceRule[T]) Validate(value any) []ValidationError {
	if value == nil {
		for _, rl := range r.rules {
			if rl.isNotNil {
				return []ValidationError{{Message: rl.message}}
			}
		}
		return nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice {
		return []ValidationError{{Message: msgArray}}
	}

	var slice []T
	if s, ok := value.([]T); ok {
		slice = s
	} else {
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			if typed, ok := elem.(T); ok {
				slice = append(slice, typed)
			} else {
				return []ValidationError{{Message: fmt.Sprintf("element [%d] has invalid type", i)}}
			}
		}
	}

	errs := executeRules(r.rules, slice, "", false)
	errs = append(errs, r.eachValueErrors(slice, "")...)

	return errs
}

// eachValueErrors validates each element with the EachValue rule, prefixing
// element field names with the slice index so nested errors keep their full
// path (e.g. "tags[0].name" rather than just "tags[0]").
func (r *SliceRule[T]) eachValueErrors(slice []T, fieldName string) []ValidationError {
	if r.eachRule == nil {
		return nil
	}

	var errs []ValidationError
	for i := range slice {
		index := fmt.Sprintf("%s[%d]", fieldName, i)
		for _, e := range r.eachRule.Validate(slice[i]) {
			errs = append(errs, ValidationError{
				Field:   joinFieldName(index, e.Field),
				Message: e.Message,
			})
		}
	}

	return errs
}

// Length rules

func (r *SliceRule[T]) NotNil() *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *SliceRule[T]) NotEmpty() *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{
		validate: func(val []T) bool { return len(val) > 0 },
		message:  msgNotEmpty,
	})
	return r
}

func (r *SliceRule[T]) Min(n int) *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{
		validate: func(val []T) bool { return len(val) >= n },
		message:  fmt.Sprintf("must have at least %d %s", n, pluralize(n, "item")),
	})
	return r
}

func (r *SliceRule[T]) Max(n int) *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{
		validate: func(val []T) bool { return len(val) <= n },
		message:  fmt.Sprintf("must have at most %d %s", n, pluralize(n, "item")),
	})
	return r
}

func (r *SliceRule[T]) Length(min, max int) *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{
		validate: func(val []T) bool { return len(val) >= min && len(val) <= max },
		message:  fmt.Sprintf(msgBetweenItems, min, max),
	})
	return r
}

// Iteration rules

func (r *SliceRule[T]) Each(fn func(item *T, sv *Validator)) *SliceRule[T] {
	r.eachFn = fn
	return r
}

func (r *SliceRule[T]) EachValue(rule Rule) *SliceRule[T] {
	r.eachRule = rule
	return r
}

// Custom rules

func (r *SliceRule[T]) Must(fn func([]T) bool) *SliceRule[T] {
	r.rules = append(r.rules, rule[[]T]{validate: fn, message: msgNotValid})
	return r
}

func (r *SliceRule[T]) WithMessage(msg string) *SliceRule[T] {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *SliceRule[T]) WithName(name string) *SliceRule[T] {
	r.name = name
	return r
}
