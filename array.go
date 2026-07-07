package validator

import (
	"fmt"
	"reflect"
)

// ArrayRule validates a fixed-size array field. Go generics can't abstract
// over an array's length, so it works by reflection and takes the field as a
// pointer. Use Slice for slices, which keeps the element type.
type ArrayRule struct {
	v        *Validator
	fieldPtr any
	name     string
	rules    []rule[int] // run against the array's length
	eachRule Rule

	// msgTarget and eachMark route WithMessage to whichever check was added last.
	eachValMsg string
	msgTarget  *string
	eachMark   int
}

// Array validates an array field by reflection; it also accepts a pointer to
// a slice. Pass the field address: without it the value is a copy and the
// name falls back to "unknown".
func Array(v *Validator, field any) *ArrayRule {
	r := &ArrayRule{v: v, fieldPtr: field}
	v.addRule(r)
	return r
}

func (r *ArrayRule) execute(prefix string) []ValidationError {
	fieldName := r.resolveName(prefix)

	rv := reflect.ValueOf(r.fieldPtr)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return executeRules(r.rules, 0, fieldName, true)
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return []ValidationError{{Field: fieldName, Message: msgArray}}
	}

	errs := containerErrors(r.rules, rv.Len(), rv.Kind() == reflect.Slice && rv.IsNil(), fieldName)

	if r.eachRule != nil {
		for i := 0; i < rv.Len(); i++ {
			index := fmt.Sprintf("%s[%d]", fieldName, i)
			for _, e := range r.eachRule.Validate(rv.Index(i).Interface()) {
				errs = append(errs, ValidationError{
					Field:   joinFieldName(index, e.Field),
					Message: orMsg(r.eachValMsg, e.Message),
				})
			}
		}
	}

	return errs
}

func (r *ArrayRule) resolveName(prefix string) string {
	if r.name != "" {
		return joinPrefix(prefix, r.name)
	}
	return resolveFieldNameForRule(r.v, r.fieldPtr, prefix)
}

// Presence rules

// NotNil fails when the field pointer is nil (e.g. a nil *[3]int) or points at
// a nil slice.
func (r *ArrayRule) NotNil() *ArrayRule {
	r.rules = append(r.rules, rule[int]{isNotNil: true, message: msgNotNil})
	return r
}

// Length rules

// NotEmpty fails when the array has no elements.
func (r *ArrayRule) NotEmpty() *ArrayRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(n int) bool { return n > 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Min requires at least n elements.
func (r *ArrayRule) Min(n int) *ArrayRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(length int) bool { return length >= n },
		message:  msgMinItems(n),
	})
	return r
}

// Max allows at most n elements.
func (r *ArrayRule) Max(n int) *ArrayRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(length int) bool { return length <= n },
		message:  msgMaxItems(n),
	})
	return r
}

// Length requires between min and max elements.
func (r *ArrayRule) Length(min, max int) *ArrayRule {
	r.rules = append(r.rules, rule[int]{
		validate: func(length int) bool { return length >= min && length <= max },
		message:  msgItemsBetween(min, max),
	})
	return r
}

// EachValue validates every element against the given rule. It panics if
// rule is nil.
func (r *ArrayRule) EachValue(rule Rule) *ArrayRule {
	if isNilRule(rule) {
		panic("validator: EachValue requires a non-nil rule")
	}
	r.eachRule = rule
	r.eachValMsg = ""
	r.msgTarget = &r.eachValMsg
	r.eachMark = len(r.rules)
	return r
}

// WithMessage replaces the previous rule's error message. Directly after
// EachValue it instead replaces the message of every error it produces.
func (r *ArrayRule) WithMessage(msg string) *ArrayRule {
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
func (r *ArrayRule) WithName(name string) *ArrayRule {
	r.name = name
	return r
}
