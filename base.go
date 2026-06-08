package validator

import (
	"reflect"
)

const (
	msgNotNil       = "must not be nil"
	msgNotEmpty     = "must not be empty"
	msgNotValid     = "is not valid"
	msgPositive     = "must be positive"
	msgNegative     = "must be negative"
	msgNonnegative  = "must be non-negative"
	msgNonpositive  = "must be non-positive"
	msgArray        = "must be an array"
	msgObject       = "must be an object"
	msgBetweenItems = "must have between %d and %d items"
)

type rule[T any] struct {
	validate func(val T) bool
	message  string
	isNotNil bool
}

// fieldBase holds the shared state and execute logic embedded by every typed
// rule; the per-type files declare only the type-specific fluent methods.
type fieldBase[T any] struct {
	v        *Validator
	fieldPtr *T
	name     string // explicit name set via WithName; overrides reflection
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

// resolveName returns the field's error name: an explicit WithName if set,
// otherwise the name resolved from the struct by reflection. Reflection can't
// resolve a nil pointer (no address), so WithName is required to name one.
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

// notNilErrors returns an error for every NotNil rule, used when a value is nil
// but the rules still need to run (e.g. a nil slice behind a non-nil pointer,
// which fails NotNil yet should still be length-checked as empty).
func notNilErrors[T any](rules []rule[T], fieldName string) []ValidationError {
	var errs []ValidationError
	for _, r := range rules {
		if r.isNotNil {
			errs = append(errs, ValidationError{Field: fieldName, Message: r.message})
		}
	}
	return errs
}

func validateStandalone[T any](rules []rule[T], value any, typeErrMsg string, coerce func(any) (T, bool)) []ValidationError {
	if value == nil {
		for _, r := range rules {
			if r.isNotNil {
				return []ValidationError{{Message: r.message}}
			}
		}
		return nil
	}

	val, ok := coerce(value)
	if !ok {
		return []ValidationError{{Message: typeErrMsg}}
	}

	var errs []ValidationError
	for _, r := range rules {
		if r.isNotNil {
			continue
		}
		if !r.validate(val) {
			errs = append(errs, ValidationError{Message: r.message})
		}
	}

	return errs
}

func resolveFieldNameForRule(v *Validator, fieldPtr any, prefix string) string {
	rv := reflect.ValueOf(fieldPtr)

	// A nil pointer has no address to match against a struct field, so its name
	// can't be resolved here — fall back to the prefix or "unknown".
	if v == nil || !rv.IsValid() || (rv.Kind() == reflect.Pointer && rv.IsNil()) {
		if prefix != "" {
			return prefix
		}
		return "unknown"
	}

	name := resolveFieldName(v.structPtr, rv.Pointer())
	if name == "" {
		name = "unknown"
	}

	return joinPrefix(prefix, name)
}

// joinPrefix qualifies a field name with its dot-notation prefix, if any.
func joinPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// pluralize returns unit, suffixed with "s" unless n is 1.
func pluralize(n int, unit string) string {
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
