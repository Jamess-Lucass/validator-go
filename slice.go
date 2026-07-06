package validator

import (
	"fmt"
	"reflect"
)

// SliceRule validates a slice field. Named slice types (type Tags []string)
// work too.
type SliceRule[S ~[]T, T any] struct {
	fieldBase[S]
	eachFn   func(item *T, sv *Validator)
	eachRule Rule

	// msgTarget and eachMark route WithMessage to whichever check was added last.
	eachFnMsg  string
	eachValMsg string
	msgTarget  *string
	eachMark   int
}

// Slice validates a slice field. Like every rule it is a free function, since
// Go methods cannot have their own type parameters. Pass the field address:
//
//	validator.Slice(v, &s.Tags).NotEmpty().EachValue(rule.String().NotEmpty())
func Slice[S ~[]T, T any](v *Validator, field *S) *SliceRule[S, T] {
	r := &SliceRule[S, T]{fieldBase: fieldBase[S]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func (r *SliceRule[S, T]) execute(prefix string) []ValidationError {
	fieldName := r.resolveName(prefix)

	if r.fieldPtr == nil {
		return executeRules(r.rules, nil, fieldName, true)
	}

	val := *r.fieldPtr

	errs := containerErrors(r.rules, val, val == nil, fieldName)

	if r.eachFn != nil {
		for i := range val {
			item := &val[i]
			sv := &Validator{structPtr: item, prefix: fmt.Sprintf("%s[%d]", fieldName, i)}
			r.eachFn(item, sv)
			for _, e := range sv.Validate().Errors {
				e.Message = orMsg(r.eachFnMsg, e.Message)
				errs = append(errs, e)
			}
		}
	}

	errs = append(errs, r.eachValueErrors(val, fieldName)...)

	return errs
}

// Validate implements Rule for standalone use. Each callbacks only run in
// struct mode, where there are field pointers to resolve names against.
func (r *SliceRule[S, T]) Validate(value any) []ValidationError {
	value = unwrapPointers(value)

	if value == nil {
		return executeRules(r.rules, nil, "", true)
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return []ValidationError{{Message: msgArray}}
	}

	var slice []T
	if s, ok := value.([]T); ok {
		slice = s
	} else {
		var typeErrs []ValidationError
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			typed, ok := elem.(T)
			if !ok {
				typeErrs = append(typeErrs, ValidationError{
					Field:   fmt.Sprintf("[%d]", i),
					Message: "is of the wrong type",
				})
				continue
			}
			slice = append(slice, typed)
		}
		// One wrong typed element means the slice can't be checked as a whole.
		if len(typeErrs) > 0 {
			return typeErrs
		}
	}

	errs := containerErrors(r.rules, S(slice), rv.Kind() == reflect.Slice && rv.IsNil(), "")
	errs = append(errs, r.eachValueErrors(S(slice), "")...)

	return errs
}

func (r *SliceRule[S, T]) eachValueErrors(slice S, fieldName string) []ValidationError {
	if r.eachRule == nil {
		return nil
	}

	var errs []ValidationError
	for i := range slice {
		index := fmt.Sprintf("%s[%d]", fieldName, i)
		for _, e := range r.eachRule.Validate(slice[i]) {
			errs = append(errs, ValidationError{
				Field:   joinFieldName(index, e.Field),
				Message: orMsg(r.eachValMsg, e.Message),
			})
		}
	}

	return errs
}

// Presence rules

// NotNil fails on a nil slice, which is distinct from an empty one.
func (r *SliceRule[S, T]) NotNil() *SliceRule[S, T] {
	r.rules = append(r.rules, rule[S]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails when the slice has no elements.
func (r *SliceRule[S, T]) NotEmpty() *SliceRule[S, T] {
	r.rules = append(r.rules, rule[S]{
		validate: func(val S) bool { return len(val) > 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Length rules

// Min requires at least n elements.
func (r *SliceRule[S, T]) Min(n int) *SliceRule[S, T] {
	r.rules = append(r.rules, rule[S]{
		validate: func(val S) bool { return len(val) >= n },
		message:  msgMinItems(n),
	})
	return r
}

// Max allows at most n elements.
func (r *SliceRule[S, T]) Max(n int) *SliceRule[S, T] {
	r.rules = append(r.rules, rule[S]{
		validate: func(val S) bool { return len(val) <= n },
		message:  msgMaxItems(n),
	})
	return r
}

// Length requires between min and max elements.
func (r *SliceRule[S, T]) Length(min, max int) *SliceRule[S, T] {
	r.rules = append(r.rules, rule[S]{
		validate: func(val S) bool { return len(val) >= min && len(val) <= max },
		message:  msgItemsBetween(min, max),
	})
	return r
}

// Iteration rules

// Each runs the callback for every element with a sub-validator; errors are
// prefixed with the index, like "items[2].name". It panics if fn is nil.
func (r *SliceRule[S, T]) Each(fn func(item *T, sv *Validator)) *SliceRule[S, T] {
	if fn == nil {
		panic("validator: Each requires a non-nil callback")
	}
	r.eachFn = fn
	r.eachFnMsg = ""
	r.msgTarget = &r.eachFnMsg
	r.eachMark = len(r.rules)
	return r
}

// EachValue checks every element against the rule, naming errors by index.
// It panics if rule is nil.
func (r *SliceRule[S, T]) EachValue(rule Rule) *SliceRule[S, T] {
	if isNilRule(rule) {
		panic("validator: EachValue requires a non-nil rule")
	}
	r.eachRule = rule
	r.eachValMsg = ""
	r.msgTarget = &r.eachValMsg
	r.eachMark = len(r.rules)
	return r
}

// Custom rules

// Must runs a custom check against the whole slice. It panics if fn is nil.
func (r *SliceRule[S, T]) Must(fn func(S) bool) *SliceRule[S, T] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[S]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message. Directly after Each
// or EachValue it instead replaces the message of every error they produce.
func (r *SliceRule[S, T]) WithMessage(msg string) *SliceRule[S, T] {
	if r.msgTarget != nil && len(r.rules) == r.eachMark {
		*r.msgTarget = msg
		return r
	}
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *SliceRule[S, T]) WithName(name string) *SliceRule[S, T] {
	r.name = name
	return r
}
