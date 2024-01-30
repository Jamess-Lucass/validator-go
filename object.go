package schema

import (
	"fmt"
	"reflect"
)

type ObjectSchema struct {
	value map[string]ISchema
}

var _ ISchema = (*ObjectSchema)(nil)

func Object(obj map[string]ISchema) *ObjectSchema {
	return &ObjectSchema{value: obj}
}

func (s *ObjectSchema) Parse(value any) *ValidationResult {
	t := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	if t.Kind() != reflect.Struct {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected struct, got %T", value)}}}
	}

	res := &ValidationResult{}

	for key, schema := range s.value {
		if _, ok := t.FieldByName(key); !ok {
			err := ValidationError{
				Path:    key,
				Message: "Required",
			}

			res.Errors = append(res.Errors, err)
			continue
		}

		str := val.FieldByName(key).Interface()

		result := schema.Parse(str)
		if !result.IsValid() {
			for _, err := range result.Errors {
				newError := ValidationError{
					Path:    key,
					Message: err.Message,
				}

				res.Errors = append(res.Errors, newError)
			}
		}
	}

	return res
}
