// Package validator provides FluentValidation-style validation for Go structs
// and unstructured data: build a Validator with New, declare rules on fields,
// then call Validate. Every rule attaches the same way, as a free function
// taking the validator and a field pointer:
//
//	validator.String(v, &user.Name).NotEmpty().Min(2)
//	validator.Number(v, &user.Age).Gte(0)
//
// Standalone rule values for unstructured data live in the rule subpackage.
//
// A Validator is built per value and is not safe for concurrent use.
// Standalone rule values are read-only once built, so they can be shared
// across goroutines.
package validator

import (
	"errors"
	"reflect"
	"strings"
)

// ValidationError is a single failed check: which field, and why.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error formats the error as "field: message".
func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

// ValidationResult holds every error collected by a Validate call.
type ValidationResult struct {
	Errors []ValidationError `json:"errors"`
}

// IsValid reports whether no errors were found.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// Rule is a validation check that runs against a standalone value. The
// constructors in the rule subpackage all return implementations.
type Rule interface {
	Validate(value any) []ValidationError
}

type fieldRule interface {
	execute(prefix string) []ValidationError
}

// Validator collects rules for one struct and runs them on Validate.
type Validator struct {
	structPtr any
	rules     []fieldRule
	prefix    string
}

// New returns a Validator for structPtr, which must be a non-nil pointer to a
// struct.
func New(structPtr any) (*Validator, error) {
	val := reflect.ValueOf(structPtr)
	if !val.IsValid() || val.Kind() != reflect.Pointer || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return nil, errors.New("validator: New requires a non-nil pointer to a struct")
	}

	return &Validator{structPtr: structPtr}, nil
}

// Validate runs every rule against the struct's current values.
func (v *Validator) Validate() *ValidationResult {
	result := &ValidationResult{}

	for _, r := range v.rules {
		errs := r.execute(v.prefix)
		result.Errors = append(result.Errors, errs...)
	}

	return result
}

func (v *Validator) addRule(r fieldRule) {
	v.rules = append(v.rules, r)
}

func resolveFieldName(structPtr any, fieldPtr uintptr) string {
	if structPtr == nil {
		return ""
	}

	// Each over []*T or map[K]*V hands over a pointer to a pointer, so unwrap
	// every level.
	val := reflect.ValueOf(structPtr)
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ""
	}

	return resolveFieldNameFromValue(val, fieldPtr, map[uintptr]struct{}{})
}

// resolveFieldNameFromValue finds the field at fieldPtr. Leaves are checked
// before recursing, since a nested struct at offset 0 shares its address with
// its first field. visited stops cyclic structures recursing forever.
func resolveFieldNameFromValue(val reflect.Value, fieldPtr uintptr, visited map[uintptr]struct{}) string {
	typ := val.Type()
	parentPkg := typ.PkgPath()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if recursesInto(field, fieldType, parentPkg) {
			continue
		}

		if field.Kind() == reflect.Pointer {
			if !field.IsNil() && field.Pointer() == fieldPtr {
				return getFieldName(fieldType)
			}
			continue
		}

		if field.Addr().Pointer() == fieldPtr {
			return getFieldName(fieldType)
		}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !recursesInto(field, fieldType, parentPkg) {
			continue
		}

		nested := field
		if field.Kind() == reflect.Pointer {
			ptr := field.Pointer()
			if _, seen := visited[ptr]; seen {
				continue
			}
			visited[ptr] = struct{}{}
			nested = field.Elem()
		}

		if name := nestedName(fieldType, resolveFieldNameFromValue(nested, fieldPtr, visited)); name != "" {
			return name
		}
	}

	return ""
}

func recursesInto(field reflect.Value, fieldType reflect.StructField, parentPkg string) bool {
	switch field.Kind() {
	case reflect.Struct:
		return recursableStruct(fieldType.Anonymous, fieldType.Type, parentPkg)
	case reflect.Pointer:
		if field.IsNil() || field.Elem().Kind() != reflect.Struct {
			return false
		}
		return recursableStruct(fieldType.Anonymous, fieldType.Type.Elem(), parentPkg)
	default:
		return false
	}
}

// Embedded fields promote without a prefix, unless a json tag nests them as
// encoding/json would.
func nestedName(parent reflect.StructField, nested string) string {
	if nested == "" {
		return ""
	}

	if parent.Anonymous {
		if name, ok := jsonTagName(parent); ok {
			return name + "." + nested
		}
		return nested
	}

	return getFieldName(parent) + "." + nested
}

// Structs with no exported fields (time.Time) stay leaves: rules point at the
// struct itself, so descending would attribute the match to the wrong field.
func recursableStruct(anonymous bool, t reflect.Type, parentPkg string) bool {
	return anonymous || t.PkgPath() == parentPkg || hasExportedFields(t)
}

func hasExportedFields(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			return true
		}
	}
	return false
}

func getFieldName(f reflect.StructField) string {
	if name, ok := jsonTagName(f); ok {
		return name
	}

	return f.Name
}

func jsonTagName(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return "", false
	}

	name, _, _ := strings.Cut(tag, ",")
	return name, name != ""
}
