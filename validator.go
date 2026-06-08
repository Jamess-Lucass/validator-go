// Package validator provides FluentValidation-style validation for Go structs
// and unstructured maps: build a Validator with New, declare rules on fields,
// and call Validate to collect any errors.
package validator

import (
	"errors"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

type Rule interface {
	Validate(value any) []ValidationError
}

type fieldRule interface {
	execute(prefix string) []ValidationError
}

type Validator struct {
	structPtr any
	rules     []fieldRule
	prefix    string
}

// New returns a Validator for the given struct pointer. It returns an error if
// structPtr is not a non-nil pointer to a struct.
func New(structPtr any) (*Validator, error) {
	val := reflect.ValueOf(structPtr)
	if !val.IsValid() || val.Kind() != reflect.Pointer || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return nil, errors.New("validator: New requires a non-nil pointer to a struct")
	}

	return &Validator{structPtr: structPtr}, nil
}

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

	val := reflect.ValueOf(structPtr)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ""
	}

	return resolveFieldNameFromValue(val, fieldPtr)
}

func resolveFieldNameFromValue(val reflect.Value, fieldPtr uintptr) string {
	typ := val.Type()
	parentPkg := typ.PkgPath()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldAddr := field.Addr().Pointer()

		if field.Kind() == reflect.Struct && isRecursableStruct(fieldType, parentPkg) {
			if name := nestedName(fieldType, resolveFieldNameFromValue(field, fieldPtr)); name != "" {
				return name
			}
		}

		if field.Kind() == reflect.Pointer && !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			elemPkg := fieldType.Type.Elem().PkgPath()
			if fieldType.Anonymous || elemPkg == parentPkg {
				if name := nestedName(fieldType, resolveFieldNameFromValue(field.Elem(), fieldPtr)); name != "" {
					return name
				}
			}
		}

		if fieldAddr == fieldPtr {
			return getFieldName(fieldType)
		}

		if field.Kind() == reflect.Pointer && !field.IsNil() {
			if field.Pointer() == fieldPtr {
				return getFieldName(fieldType)
			}
		}
	}

	return ""
}

// nestedName qualifies a resolved nested field name with its parent field name,
// except for anonymous (embedded) fields which are promoted without a prefix.
// It returns "" when nothing matched in the nested struct.
func nestedName(parent reflect.StructField, nested string) string {
	if nested == "" {
		return ""
	}

	if parent.Anonymous {
		return nested
	}

	return getFieldName(parent) + "." + nested
}

// Only recurse into anonymous (embedded) fields or structs defined in the same
// package. External types like time.Time or uuid.UUID have unexported fields
// that would panic during reflection.
func isRecursableStruct(f reflect.StructField, parentPkg string) bool {
	if f.Anonymous {
		return true
	}
	return f.Type.PkgPath() == parentPkg
}

func getFieldName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag != "" && tag != "-" {
		name := strings.Split(tag, ",")[0]
		if name != "" {
			return name
		}
	}

	return f.Name
}
