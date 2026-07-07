package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	msgNotNil      = "must not be nil"
	msgNotEmpty    = "must not be empty"
	msgNotValid    = "is not valid"
	msgPositive    = "must be positive"
	msgNegative    = "must be negative"
	msgNonnegative = "must be non-negative"
	msgNonpositive = "must be non-positive"
	msgArray       = "must be an array"
	msgObject      = "must be an object"
)

func msgMinItems(n int) string {
	return fmt.Sprintf("must have at least %d %s", n, pluralise(n, "item"))
}

func msgMaxItems(n int) string {
	return fmt.Sprintf("must have at most %d %s", n, pluralise(n, "item"))
}

func msgItemsBetween(min, max int) string {
	return fmt.Sprintf("must have between %d and %d items", min, max)
}

// Strings are quoted so empty or spaced values stay visible; an empty list
// falls back to the generic message.
func msgOneOf[T any](values []T, quoted bool) string {
	if len(values) == 0 {
		return msgNotValid
	}

	parts := make([]string, len(values))
	for i, v := range values {
		if quoted {
			parts[i] = fmt.Sprintf("%q", fmt.Sprint(v))
		} else {
			parts[i] = fmt.Sprint(v)
		}
	}
	return "must be one of " + strings.Join(parts, ", ")
}

type rule[T any] struct {
	validate func(val T) bool
	message  string
	isNotNil bool
}

// fieldBase carries the state and execute logic shared by every typed rule.
type fieldBase[T any] struct {
	v        *Validator
	fieldPtr *T
	name     string
	rules    []rule[T]
}

func (b *fieldBase[T]) execute(prefix string) []ValidationError {
	fieldName := b.resolveName(prefix)

	isNil := b.fieldPtr == nil
	var val T
	if !isNil {
		val = *b.fieldPtr
	}

	return executeRules(b.rules, val, fieldName, isNil)
}

func (b *fieldBase[T]) resolveName(prefix string) string {
	if b.name != "" {
		return joinPrefix(prefix, b.name)
	}
	return resolveFieldNameForRule(b.v, b.fieldPtr, prefix)
}

func executeRules[T any](rules []rule[T], val T, fieldName string, isNil bool) []ValidationError {
	if isNil {
		for _, r := range rules {
			if r.isNotNil {
				return []ValidationError{{Field: fieldName, Message: r.message}}
			}
		}
		return nil
	}

	var errs []ValidationError
	for _, r := range rules {
		if r.isNotNil {
			continue
		}
		if !r.validate(val) {
			errs = append(errs, ValidationError{Field: fieldName, Message: r.message})
		}
	}

	return errs
}

// containerErrors runs rules against a container that exists but may itself
// be nil: a nil slice or map fails NotNil yet still length checks as empty.
func containerErrors[T any](rules []rule[T], val T, valIsNil bool, fieldName string) []ValidationError {
	var errs []ValidationError
	for _, r := range rules {
		if r.isNotNil {
			if valIsNil {
				errs = append(errs, ValidationError{Field: fieldName, Message: r.message})
			}
			continue
		}
		if !r.validate(val) {
			errs = append(errs, ValidationError{Field: fieldName, Message: r.message})
		}
	}
	return errs
}

func orMsg(override, msg string) string {
	if override != "" {
		return override
	}
	return msg
}

func validateStandalone[T any](rules []rule[T], value any, typeErrMsg string, coerce func(any) (T, bool)) []ValidationError {
	value = unwrapPointers(value)

	if value == nil {
		var zero T
		return executeRules(rules, zero, "", true)
	}

	val, ok := coerce(value)
	if !ok {
		return []ValidationError{{Message: typeErrMsg}}
	}

	return executeRules(rules, val, "", false)
}

// isNilRule catches both nil and a typed nil pointer boxed in the interface.
func isNilRule(r Rule) bool {
	if r == nil {
		return true
	}

	rv := reflect.ValueOf(r)
	return rv.Kind() == reflect.Pointer && rv.IsNil()
}

// unwrapPointers dereferences any level of pointer; a typed nil becomes plain nil.
func unwrapPointers(value any) any {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer {
		return value
	}

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	return rv.Interface()
}

func resolveFieldNameForRule(v *Validator, fieldPtr any, prefix string) string {
	rv := reflect.ValueOf(fieldPtr)

	if v == nil || rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fallbackName(prefix)
	}

	name := resolveFieldName(v.structPtr, rv.Pointer())
	if name == "" {
		return fallbackName(prefix)
	}

	return joinPrefix(prefix, name)
}

func fallbackName(prefix string) string {
	if prefix != "" {
		return prefix
	}
	return "unknown"
}

func joinPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

func pluralise(n int, unit string) string {
	if n == 1 {
		return unit
	}
	return unit + "s"
}

func joinFieldName(parent, child string) string {
	if child == "" {
		return parent
	}

	if parent == "" {
		return child
	}

	if child[0] == '[' {
		return parent + child
	}

	return parent + "." + child
}
